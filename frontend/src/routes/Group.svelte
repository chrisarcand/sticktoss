<script>
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import { playersAPI, groupsAPI } from '../lib/api';
  import { skillLevels } from '../lib/store';

  export let id;

  let group = null;
  let allPlayers = [];
  let loading = true;
  let showAddPlayerModal = false;
  let showTeams = false;
  let teams = [];
  let numTeams = 2;
  let showWeights = false;
  let selectedPlayers = new Set();
  let lockedGroups = []; // Array of arrays of player IDs

  onMount(async () => {
    if (!localStorage.getItem('token')) {
      navigate('/login');
      return;
    }

    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      [group, allPlayers] = await Promise.all([
        groupsAPI.get(id),
        playersAPI.getAll()
      ]);
    } catch (err) {
      console.error('Failed to load data:', err);
      alert('Failed to load group');
      navigate('/');
    } finally {
      loading = false;
    }
  }

  async function addPlayerToGroup(playerId) {
    try {
      await groupsAPI.addPlayer(id, playerId);
      await loadData();
      showAddPlayerModal = false;
    } catch (err) {
      alert(err.message);
    }
  }

  async function removePlayerFromGroup(playerId) {
    if (confirm('Remove this player from the group?')) {
      try {
        await groupsAPI.removePlayer(id, playerId);
        await loadData();
      } catch (err) {
        alert(err.message);
      }
    }
  }

  async function generateTeams() {
    if (!group.players || group.players.length < numTeams) {
      alert(`You need at least ${numTeams} players to generate ${numTeams} teams.`);
      return;
    }

    // Convert locked groups to arrays of player IDs
    const lockedPlayerIds = lockedGroups.filter(g => g.length > 0);

    try {
      const result = await groupsAPI.generateTeams(id, numTeams, lockedPlayerIds);
      teams = result.teams;
      showTeams = true;
      showWeights = false;
      selectedPlayers.clear();
      lockedGroups = [];
    } catch (err) {
      alert(err.message);
    }
  }

  function togglePlayerSelection(playerId) {
    if (selectedPlayers.has(playerId)) {
      selectedPlayers.delete(playerId);
    } else {
      selectedPlayers.add(playerId);
    }
    selectedPlayers = selectedPlayers; // Trigger reactivity
  }

  function lockSelectedPlayers() {
    if (selectedPlayers.size === 0) {
      alert('Select at least one player to lock together');
      return;
    }

    lockedGroups = [...lockedGroups, Array.from(selectedPlayers)];
    selectedPlayers.clear();
    selectedPlayers = selectedPlayers;
  }

  function removeLockGroup(index) {
    lockedGroups = lockedGroups.filter((_, i) => i !== index);
  }

  function getPlayerName(playerId) {
    const player = group.players.find(p => p.id === playerId);
    return player ? player.name : '';
  }

  $: availablePlayers = allPlayers.filter(p =>
    !group?.players?.some(gp => gp.id === p.id)
  );
</script>

<div class="container">
  <div class="header">
    <button class="back-btn" on:click={() => navigate('/')}>← Back</button>
    <h1>{group?.name || 'Loading...'}</h1>
  </div>

  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="content">
      <section class="section">
        <div class="section-header">
          <h2>Players in Group ({group.players?.length || 0})</h2>
          <button on:click={() => showAddPlayerModal = true}>+ Add Player</button>
        </div>

        {#if !group.players || group.players.length === 0}
          <p class="empty">No players in this group yet. Add some players to get started!</p>
        {:else}
          <div class="players-list">
            {#each group.players as player}
              <div class="player-card" class:selected={selectedPlayers.has(player.id)}>
                <div class="player-info" on:click={() => togglePlayerSelection(player.id)}>
                  <div class="checkbox">
                    {#if selectedPlayers.has(player.id)}
                      <span class="checkmark">✓</span>
                    {/if}
                  </div>
                  <div class="player-details">
                    <h3>{player.name}</h3>
                    <div class="weight-badge">
                      Level {player.skill_weight} - {skillLevels[player.skill_weight].label}
                    </div>
                  </div>
                </div>
                <button class="btn-remove" on:click={() => removePlayerFromGroup(player.id)}>Remove</button>
              </div>
            {/each}
          </div>

          <div class="lock-controls">
            <button
              class="lock-btn"
              on:click={lockSelectedPlayers}
              disabled={selectedPlayers.size === 0}
            >
              Lock Selected Players Together
            </button>
            {#if selectedPlayers.size > 0}
              <span class="selection-count">{selectedPlayers.size} selected</span>
            {/if}
          </div>

          {#if lockedGroups.length > 0}
            <div class="locked-groups">
              <h3>Locked Groups:</h3>
              {#each lockedGroups as lockGroup, index}
                <div class="locked-group">
                  <span>{lockGroup.map(id => getPlayerName(id)).join(', ')}</span>
                  <button class="btn-small-remove" on:click={() => removeLockGroup(index)}>×</button>
                </div>
              {/each}
            </div>
          {/if}
        {/if}
      </section>

      {#if group.players && group.players.length >= 2}
        <section class="section">
          <h2>Generate Teams</h2>
          <div class="generate-controls">
            <div class="form-group">
              <label>Number of Teams:</label>
              <input type="number" bind:value={numTeams} min="2" max={group.players.length} />
            </div>
            <button class="generate-btn" on:click={generateTeams}>Generate Teams</button>
          </div>
        </section>
      {:else if group.players && group.players.length === 1}
        <section class="section">
          <div class="empty-state">
            <h3>Almost there!</h3>
            <p>You need at least 2 players in this group to generate teams. Add one more player to get started.</p>
          </div>
        </section>
      {:else if group.players && group.players.length === 0}
        <section class="section">
          <div class="empty-state">
            <h3>Ready to generate teams?</h3>
            <p>Add at least 2 players to this group to start generating balanced teams.</p>
          </div>
        </section>
      {/if}

      {#if showTeams && teams.length > 0}
        <section class="section">
          <div class="section-header">
            <h2>Generated Teams</h2>
            <button class="toggle-weights" on:click={() => showWeights = !showWeights}>
              {showWeights ? 'Hide' : 'Show'} Weights
            </button>
          </div>

          <div class="teams-grid">
            {#each teams as team}
              <div class="team-card">
                <h3>Team {team.number}</h3>
                {#if showWeights}
                  <div class="total-weight">Total Weight: {team.total_weight}</div>
                {/if}
                <ul class="team-players">
                  {#each team.players as player}
                    <li>
                      {player.name}
                      {#if showWeights}
                        <span class="player-weight">({player.skill_weight})</span>
                      {/if}
                    </li>
                  {/each}
                </ul>
              </div>
            {/each}
          </div>
        </section>
      {/if}
    </div>
  {/if}
</div>

<!-- Add Player Modal -->
{#if showAddPlayerModal}
  <div class="modal-overlay" on:click={() => showAddPlayerModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Add Player to Group</h2>
      {#if availablePlayers.length === 0}
        <p class="empty">No available players. Create players from the dashboard first!</p>
        <button on:click={() => showAddPlayerModal = false}>Close</button>
      {:else}
        <div class="player-list">
          {#each availablePlayers as player}
            <div class="player-list-item" on:click={() => addPlayerToGroup(player.id)}>
              <span>{player.name}</span>
              <span class="weight-badge">Level {player.skill_weight}</span>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .header {
    margin-bottom: 30px;
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

  .loading {
    text-align: center;
    padding: 40px;
    font-size: 18px;
    color: #666;
  }

  .section {
    background: white;
    border-radius: 8px;
    padding: 30px;
    margin-bottom: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .section-header h2 {
    font-size: 24px;
    color: #333;
    margin: 0;
  }

  .section-header button {
    padding: 10px 20px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .section-header button:hover {
    background-color: #45a049;
  }

  .toggle-weights {
    background-color: #9E9E9E !important;
  }

  .toggle-weights:hover {
    background-color: #757575 !important;
  }

  .empty {
    text-align: center;
    padding: 40px;
    color: #999;
    font-style: italic;
  }

  .empty-state {
    text-align: center;
    padding: 40px;
    background-color: #f9f9f9;
    border-radius: 8px;
    border: 2px dashed #e0e0e0;
  }

  .empty-state h3 {
    margin: 0 0 10px 0;
    color: #333;
    font-size: 20px;
  }

  .empty-state p {
    margin: 0;
    color: #666;
    font-size: 16px;
  }

  .players-list {
    display: grid;
    gap: 15px;
  }

  .player-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    transition: all 0.2s;
  }

  .player-card:hover {
    border-color: #bdbdbd;
  }

  .player-card.selected {
    border-color: #4CAF50;
    background-color: #f1f8f4;
  }

  .player-info {
    display: flex;
    align-items: center;
    gap: 15px;
    flex: 1;
    cursor: pointer;
  }

  .checkbox {
    width: 24px;
    height: 24px;
    border: 2px solid #bdbdbd;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .player-card.selected .checkbox {
    background-color: #4CAF50;
    border-color: #4CAF50;
  }

  .checkmark {
    color: white;
    font-weight: bold;
  }

  .player-details h3 {
    margin: 0 0 5px 0;
    font-size: 16px;
    color: #333;
  }

  .weight-badge {
    display: inline-block;
    background-color: #e3f2fd;
    color: #1976d2;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }

  .btn-remove {
    padding: 6px 12px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .btn-remove:hover {
    background-color: #da190b;
  }

  .lock-controls {
    margin-top: 20px;
    display: flex;
    align-items: center;
    gap: 15px;
  }

  .lock-btn {
    padding: 10px 20px;
    background-color: #FF9800;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .lock-btn:hover:not(:disabled) {
    background-color: #F57C00;
  }

  .lock-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .selection-count {
    color: #666;
    font-size: 14px;
  }

  .locked-groups {
    margin-top: 20px;
    padding: 15px;
    background-color: #fff3e0;
    border-radius: 4px;
  }

  .locked-groups h3 {
    margin: 0 0 10px 0;
    font-size: 16px;
    color: #e65100;
  }

  .locked-group {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background-color: white;
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .btn-small-remove {
    background: none;
    border: none;
    color: #f44336;
    font-size: 20px;
    cursor: pointer;
    padding: 0 8px;
  }

  .generate-controls {
    display: flex;
    align-items: flex-end;
    gap: 20px;
  }

  .form-group {
    flex: 0 0 200px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    color: #333;
    font-weight: 500;
  }

  .form-group input {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
  }

  .generate-btn {
    padding: 10px 30px;
    background-color: #2196F3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    font-weight: 500;
  }

  .generate-btn:hover {
    background-color: #0b7dda;
  }

  .teams-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
  }

  .team-card {
    background-color: #f5f5f5;
    padding: 20px;
    border-radius: 8px;
    border: 2px solid #e0e0e0;
  }

  .team-card h3 {
    margin: 0 0 10px 0;
    font-size: 20px;
    color: #333;
  }

  .total-weight {
    font-weight: bold;
    color: #666;
    margin-bottom: 15px;
    padding: 8px;
    background-color: white;
    border-radius: 4px;
  }

  .team-players {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .team-players li {
    padding: 8px 0;
    border-bottom: 1px solid #e0e0e0;
  }

  .team-players li:last-child {
    border-bottom: none;
  }

  .player-weight {
    color: #666;
    font-size: 14px;
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

  .player-list {
    max-height: 400px;
    overflow-y: auto;
  }

  .player-list-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    margin-bottom: 10px;
    cursor: pointer;
    transition: all 0.2s;
  }

  .player-list-item:hover {
    background-color: #f5f5f5;
    border-color: #4CAF50;
  }

  @media (max-width: 768px) {
    .teams-grid {
      grid-template-columns: 1fr;
    }

    .generate-controls {
      flex-direction: column;
      align-items: stretch;
    }

    .form-group {
      flex: 1;
    }

    .section-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 10px;
    }

    .section-header button {
      width: 100%;
    }
  }
</style>
