<script>
  import { navigate } from 'svelte-routing';
  import { onMount } from 'svelte';
  import { groupsAPI, playersAPI, gameAPI } from '../lib/api';
  import { skillLevels } from '../lib/store';

  export let id = undefined; // For /group/:id/teams route
  export let shareId = undefined; // For /game/:shareId route

  let teams = [];
  let groupName = '';
  let groupId = '';
  let showWeights = false;
  let numTeams = 2;
  let lockedPlayerIds = [];
  let separatedPlayerIds = [];
  let loading = false;
  let showEditPlayerModal = false;
  let editingPlayer = null;
  let useJerseyColors = false;
  let isPublicMode = false; // Whether this is a public shared game
  let currentShareId = ''; // Store the share ID for sharing
  let linkCopied = false; // Track if link was just copied
  let logoUrl = ''; // Logo URL if available

  onMount(async () => {
    // Determine mode: public (shareId) or authenticated (id)
    if (shareId) {
      // Public mode: load game from API
      isPublicMode = true;
      currentShareId = shareId;
      await loadPublicGame();
    } else {
      // Authenticated mode: load from sessionStorage
      isPublicMode = false;
      loadFromSessionStorage();
    }
  });

  async function loadPublicGame() {
    loading = true;
    try {
      const data = await gameAPI.get(shareId);
      teams = data.teams || [];
      groupName = data.group_name || 'Unknown Group';
      numTeams = data.num_teams || 2;
      useJerseyColors = data.use_jersey_colors || false;
      currentShareId = data.share_id;

      // Load logo if available
      if (data.has_logo) {
        logoUrl = `/api/game/${shareId}/logo`;
      } else {
        logoUrl = '';
      }
    } catch (err) {
      alert('Failed to load game: ' + err.message);
      navigate('/');
    } finally {
      loading = false;
    }
  }

  function loadFromSessionStorage() {
    // Get teams from sessionStorage
    const storedData = sessionStorage.getItem('generatedTeams');
    if (storedData) {
      const data = JSON.parse(storedData);
      teams = data.teams || [];
      groupName = data.groupName || 'Unknown Group';
      groupId = data.groupId || id;
      numTeams = data.numTeams || 2;
      lockedPlayerIds = data.lockedPlayerIds || [];
      separatedPlayerIds = data.separatedPlayerIds || [];
      useJerseyColors = data.useJerseyColors || false;
      currentShareId = data.shareId || '';

      // Set logo URL if we have a groupId (will be checked via image loading)
      if (groupId) {
        logoUrl = `/api/groups/${groupId}/logo`;
      } else {
        logoUrl = '';
      }

      // Clear sessionStorage after reading
      sessionStorage.removeItem('generatedTeams');
    } else {
      // No teams data, redirect back
      navigate('/');
    }
  }

  function goBack() {
    if (isPublicMode) {
      // In public mode, don't go back (no auth)
      return;
    }
    if (groupId) {
      navigate(`/group/${groupId}`);
    } else {
      navigate('/');
    }
  }

  function copyShareLink() {
    const shareUrl = `${window.location.origin}/game/${currentShareId}`;
    navigator.clipboard.writeText(shareUrl).then(() => {
      linkCopied = true;
      setTimeout(() => {
        linkCopied = false;
      }, 2000);
    }).catch(() => {
      // Fallback for older browsers
      prompt('Copy this link:', shareUrl);
    });
  }

  async function regenerate() {
    if (!groupId) return;

    loading = true;
    try {
      const result = await groupsAPI.generateTeams(groupId, numTeams, lockedPlayerIds, separatedPlayerIds);
      teams = result.teams;
    } catch (err) {
      alert(err.message);
    } finally {
      loading = false;
    }
  }

  function openEditPlayer(player) {
    editingPlayer = { ...player };
    showEditPlayerModal = true;
  }

  async function updatePlayer() {
    try {
      await playersAPI.update(editingPlayer.id, {
        name: editingPlayer.name,
        skill_weight: editingPlayer.skill_weight
      });

      // Update player in teams array
      teams = teams.map(team => ({
        ...team,
        players: team.players.map(p =>
          p.id === editingPlayer.id ? { ...p, ...editingPlayer } : p
        ),
        // Recalculate total weight
        total_weight: team.players.reduce((sum, p) =>
          sum + (p.id === editingPlayer.id ? editingPlayer.skill_weight : p.skill_weight), 0
        )
      }));

      showEditPlayerModal = false;
      editingPlayer = null;
    } catch (err) {
      alert(err.message);
    }
  }

  function getTeamName(teamNumber) {
    if (useJerseyColors && teams.length === 2) {
      return teamNumber === 1 ? 'Team Light' : 'Team Dark';
    }
    return `Team ${teamNumber}`;
  }

  function getJerseyIcon(teamNumber) {
    if (!useJerseyColors || teams.length !== 2) return null;

    const isDark = teamNumber === 2;
    const color = isDark ? '#333333' : '#f5f5f5';
    const stroke = isDark ? '#666666' : '#cccccc';

    return `<svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path d="M16 4h3l2 2v4l-2 2h-1v8h-12v-8h-1l-2-2v-4l2-2h3c0-1.1 0.9-2 2-2h4c1.1 0 2 0.9 2 2z"
            fill="${color}"
            stroke="${stroke}"
            stroke-width="1"/>
    </svg>`;
  }
</script>

<div class="container">
  <div class="header">
    {#if !isPublicMode}
      <button class="back-btn" on:click={goBack}>← Back to Group</button>
    {/if}
    <div class="header-content">
      {#if logoUrl}
        <img src={logoUrl} alt="{groupName} logo" class="group-logo" on:error={() => logoUrl = ''} />
      {/if}
      <h1>{groupName}</h1>
    </div>
  </div>

  {#if teams.length > 0}
    <div class="teams-header">
      <h2>Team Assignments</h2>
      <div class="controls">
        {#if currentShareId}
          <button class="share-btn" on:click={copyShareLink}>
            {#if linkCopied}
              ✓ Copied!
            {:else}
              📋 Copy Share Link
            {/if}
          </button>
        {/if}
        {#if !isPublicMode}
          <button class="toggle-weights-btn" on:click={() => showWeights = !showWeights}>
            {showWeights ? 'Hide' : 'Show'} Weights
          </button>
          <button class="regenerate-btn" on:click={regenerate} disabled={loading}>
            {loading ? 'Regenerating...' : 'Regenerate Teams'}
          </button>
        {/if}
      </div>
    </div>

    <div class="teams-grid">
      {#each teams as team}
        <div class="team-card">
          <h3>
            {#if getJerseyIcon(team.number)}
              <span class="jersey-icon">{@html getJerseyIcon(team.number)}</span>
            {/if}
            {getTeamName(team.number)}
          </h3>
          {#if !isPublicMode && showWeights}
            <div class="total-weight">Total Weight: {team.total_weight}</div>
          {/if}
          <ul class="team-players">
            {#each team.players as player}
              <li>
                {#if !isPublicMode}
                  <button class="player-name-btn" on:click={() => openEditPlayer(player)}>
                    {player.name}
                  </button>
                {:else}
                  <span class="player-name">{player.name}</span>
                {/if}
                {#if !isPublicMode && showWeights}
                  <span class="player-weight">
                    Level {player.skill_weight} - {skillLevels[player.skill_weight].label}
                  </span>
                {/if}
              </li>
            {/each}
          </ul>
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty">
      <p>No teams generated yet.</p>
      <button on:click={goBack}>Go Back</button>
    </div>
  {/if}
</div>

<!-- Edit Player Modal -->
{#if showEditPlayerModal && editingPlayer}
  <div class="modal-overlay" on:click={() => showEditPlayerModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Edit Player</h2>
      <form on:submit|preventDefault={updatePlayer}>
        <div class="form-group">
          <label>Name</label>
          <input type="text" bind:value={editingPlayer.name} required />
        </div>
        <div class="form-group">
          <label>
            Skill Level
            <span class="info-icon" title={skillLevels[editingPlayer.skill_weight].description}>ℹ️</span>
          </label>
          <select bind:value={editingPlayer.skill_weight} required>
            {#each Object.entries(skillLevels) as [level, info]}
              <option value={Number(level)}>{level} - {info.label}</option>
            {/each}
          </select>
          <p class="skill-description">{skillLevels[editingPlayer.skill_weight].description}</p>
        </div>
        <div class="modal-actions">
          <button type="submit" class="btn-primary">Save Changes</button>
          <button type="button" on:click={() => showEditPlayerModal = false}>Cancel</button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    min-height: 100vh;
  }

  .header {
    margin-bottom: 30px;
  }

  .header-content {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .group-logo {
    width: 80px;
    height: 80px;
    object-fit: contain;
    border-radius: 8px;
    background-color: #f5f5f5;
    padding: 8px;
  }

  .back-btn {
    background: none;
    border: none;
    color: #2196F3;
    font-size: 16px;
    cursor: pointer;
    padding: 5px 0;
    margin-bottom: 10px;
  }

  .back-btn:hover {
    text-decoration: underline;
  }

  h1 {
    font-size: 36px;
    color: #333;
    margin: 0;
  }

  .teams-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    padding-bottom: 15px;
    border-bottom: 2px solid #e0e0e0;
  }

  .teams-header h2 {
    font-size: 28px;
    color: #333;
    margin: 0;
  }

  .controls {
    display: flex;
    gap: 10px;
  }

  .toggle-weights-btn {
    padding: 10px 20px;
    background-color: #9E9E9E;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .toggle-weights-btn:hover {
    background-color: #757575;
  }

  .regenerate-btn {
    padding: 10px 20px;
    background-color: #2196F3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .regenerate-btn:hover:not(:disabled) {
    background-color: #0b7dda;
  }

  .regenerate-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .share-btn {
    padding: 10px 20px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .share-btn:hover {
    background-color: #45a049;
  }

  .teams-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 25px;
    margin-bottom: 40px;
  }

  .team-card {
    background-color: white;
    padding: 25px;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border: 2px solid #e0e0e0;
  }

  .team-card h3 {
    margin: 0 0 15px 0;
    font-size: 24px;
    color: #333;
    padding-bottom: 10px;
    border-bottom: 2px solid #4CAF50;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .jersey-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .jersey-icon svg {
    display: block;
  }

  .total-weight {
    font-weight: bold;
    color: #666;
    margin-bottom: 15px;
    padding: 10px;
    background-color: #f5f5f5;
    border-radius: 4px;
    text-align: center;
  }

  .team-players {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .team-players li {
    padding: 12px 0;
    border-bottom: 1px solid #e0e0e0;
    font-size: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .team-players li:last-child {
    border-bottom: none;
  }

  .player-name-btn {
    background: none;
    border: none;
    padding: 0;
    font-weight: 500;
    color: #2196F3;
    cursor: pointer;
    text-align: left;
    font-size: 16px;
  }

  .player-name-btn:hover {
    text-decoration: underline;
  }

  .player-name {
    font-weight: 500;
    color: #333;
    font-size: 16px;
  }

  .player-weight {
    color: #666;
    font-size: 14px;
  }

  .footer {
    text-align: center;
    margin-top: 40px;
    padding-top: 20px;
    border-top: 2px solid #e0e0e0;
  }

  .hint {
    color: #999;
    font-style: italic;
    margin: 0;
  }

  .empty {
    text-align: center;
    padding: 60px 20px;
  }

  .empty p {
    font-size: 18px;
    color: #666;
    margin-bottom: 20px;
  }

  .empty button {
    padding: 12px 24px;
    background-color: #2196F3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }

  .empty button:hover {
    background-color: #0b7dda;
  }

  @media (max-width: 768px) {
    .teams-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 15px;
    }

    .controls {
      width: 100%;
      flex-direction: column;
    }

    .toggle-weights-btn,
    .regenerate-btn {
      width: 100%;
    }

    .teams-grid {
      grid-template-columns: 1fr;
    }
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 20px;
  }

  .modal {
    background: white;
    padding: 30px;
    border-radius: 8px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal h2 {
    margin: 0 0 20px 0;
    font-size: 24px;
    color: #333;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    color: #333;
    font-weight: 500;
  }

  .form-group input,
  .form-group select {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
  }

  .info-icon {
    cursor: help;
    margin-left: 5px;
  }

  .skill-description {
    margin: 10px 0 0 0;
    padding: 10px;
    background-color: #f5f5f5;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    font-style: italic;
  }

  .modal-actions {
    display: flex;
    gap: 10px;
    margin-top: 20px;
  }

  .modal-actions button {
    flex: 1;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }

  .btn-primary {
    background-color: #4CAF50;
    color: white;
  }

  .btn-primary:hover {
    background-color: #45a049;
  }

  .modal-actions button[type="button"] {
    background-color: #e0e0e0;
    color: #333;
  }

  .modal-actions button[type="button"]:hover {
    background-color: #d0d0d0;
  }

  @media print {
    .back-btn,
    .controls,
    .footer {
      display: none;
    }

    .team-card {
      break-inside: avoid;
    }
  }
</style>
