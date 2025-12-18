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
  let numTeams = 2;
  let selectedPlayers = new Set();
  let lockedGroups = []; // Array of arrays of player IDs
  let separatedGroups = []; // Array of arrays of player IDs that must be on different teams
  let useJerseyColors = true;
  let selectedPlayersToAdd = new Set();
  let showWeightBadges = false;
  let logoUrl = '';
  let uploadingLogo = false;

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
      // Set logo URL if group has a logo
      if (group.logo_content_type) {
        logoUrl = `/api/groups/${id}/logo`;
      } else {
        logoUrl = '';
      }
    } catch (err) {
      console.error('Failed to load data:', err);
      alert('Failed to load group');
      navigate('/');
    } finally {
      loading = false;
    }
  }

  async function handleLogoUpload(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Validate file type
    if (!file.type.match(/^image\/(png|jpeg|svg\+xml)$/)) {
      alert('Please upload a PNG, JPG, or SVG image');
      return;
    }

    // Validate file size (2MB)
    if (file.size > 2 * 1024 * 1024) {
      alert('File too large. Maximum size is 2MB');
      return;
    }

    uploadingLogo = true;
    try {
      await groupsAPI.uploadLogo(id, file);
      await loadData(); // Reload to get updated logo
    } catch (err) {
      alert(err.message);
    } finally {
      uploadingLogo = false;
    }
  }

  async function deleteLogo() {
    if (!confirm('Remove logo from this group?')) return;

    try {
      await groupsAPI.deleteLogo(id);
      await loadData();
    } catch (err) {
      alert(err.message);
    }
  }

  async function addPlayersToGroup() {
    if (selectedPlayersToAdd.size === 0) {
      alert('Please select at least one player to add');
      return;
    }

    try {
      // Add all selected players
      await Promise.all(
        Array.from(selectedPlayersToAdd).map(playerId =>
          groupsAPI.addPlayer(id, playerId)
        )
      );
      await loadData();
      selectedPlayersToAdd.clear();
      selectedPlayersToAdd = selectedPlayersToAdd;
      showAddPlayerModal = false;
    } catch (err) {
      alert(err.message);
    }
  }

  function togglePlayerToAdd(playerId) {
    if (selectedPlayersToAdd.has(playerId)) {
      selectedPlayersToAdd.delete(playerId);
    } else {
      selectedPlayersToAdd.add(playerId);
    }
    selectedPlayersToAdd = selectedPlayersToAdd;
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

  async function removeSelectedPlayers() {
    const count = selectedPlayers.size;
    if (confirm(`Remove ${count} selected player${count > 1 ? 's' : ''} from the group?`)) {
      try {
        // Remove all selected players
        await Promise.all(
          Array.from(selectedPlayers).map(playerId =>
            groupsAPI.removePlayer(id, playerId)
          )
        );
        selectedPlayers.clear();
        selectedPlayers = selectedPlayers;
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

    // Convert locked and separated groups to arrays of player IDs
    const lockedPlayerIds = lockedGroups.filter(g => g.length > 0);
    const separatedPlayerIds = separatedGroups.filter(g => g.length > 0);
    const jerseyColors = useJerseyColors && numTeams === 2;

    try {
      const result = await groupsAPI.generateTeams(id, numTeams, lockedPlayerIds, separatedPlayerIds, jerseyColors);

      // Store teams in sessionStorage and navigate
      sessionStorage.setItem('generatedTeams', JSON.stringify({
        teams: result.teams,
        shareId: result.share_id,
        groupName: group.name,
        groupId: id,
        numTeams: numTeams,
        lockedPlayerIds: lockedPlayerIds,
        separatedPlayerIds: separatedPlayerIds,
        useJerseyColors: jerseyColors
      }));

      navigate(`/group/${id}/teams`);
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

    const newLock = Array.from(selectedPlayers);

    // Check for exact duplicate
    const isDuplicate = lockedGroups.some(group =>
      group.length === newLock.length &&
      group.every(id => newLock.includes(id))
    );

    if (isDuplicate) {
      alert('This exact lock already exists');
      return;
    }

    // Merge overlapping groups (transitive closure)
    const mergedGroups = [];
    const playersToMerge = new Set(newLock);

    for (const group of lockedGroups) {
      // Check if this group overlaps with our new lock
      const hasOverlap = group.some(id => playersToMerge.has(id));

      if (hasOverlap) {
        // Merge this group into our set
        group.forEach(id => playersToMerge.add(id));
      } else {
        // Keep this group separate
        mergedGroups.push(group);
      }
    }

    // Add the merged group
    const finalMergedGroup = Array.from(playersToMerge);

    // Validate: check if merged group is too large for any team
    // (merged group must fit on a single team)
    const playersPerTeam = Math.ceil(group.players.length / numTeams);
    if (finalMergedGroup.length > playersPerTeam) {
      alert(`Cannot lock ${finalMergedGroup.length} players together. With ${numTeams} teams and ${group.players.length} players, each team can have at most ${playersPerTeam} players.`);
      return;
    }

    mergedGroups.push(finalMergedGroup);
    lockedGroups = mergedGroups;

    selectedPlayers.clear();
    selectedPlayers = selectedPlayers;
  }

  function removeLockGroup(index) {
    lockedGroups = lockedGroups.filter((_, i) => i !== index);
  }

  function separateSelectedPlayers() {
    if (selectedPlayers.size === 0) {
      alert('Select at least 2 players to keep apart');
      return;
    }

    if (selectedPlayers.size < 2) {
      alert('You need at least 2 players to separate');
      return;
    }

    if (selectedPlayers.size > numTeams) {
      alert(`Cannot separate ${selectedPlayers.size} players with only ${numTeams} teams. You can separate at most ${numTeams} players.`);
      return;
    }

    const newSeparation = Array.from(selectedPlayers);

    // Check for duplicates
    const isDuplicate = separatedGroups.some(group =>
      group.length === newSeparation.length &&
      group.every(id => newSeparation.includes(id))
    );

    if (isDuplicate) {
      alert('This exact separation already exists');
      return;
    }

    // Check for overlap with locked groups
    for (const playerId of newSeparation) {
      const isLocked = lockedGroups.some(group => group.includes(playerId));
      if (isLocked) {
        alert('Cannot separate players that are locked together');
        return;
      }
    }

    separatedGroups = [...separatedGroups, newSeparation];
    selectedPlayers.clear();
    selectedPlayers = selectedPlayers;
  }

  function removeSeparatedGroup(index) {
    separatedGroups = separatedGroups.filter((_, i) => i !== index);
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
    <div class="header-content">
      {#if logoUrl}
        <img src={logoUrl} alt="{group?.name} logo" class="group-logo" />
      {/if}
      <h1>{group?.name || 'Loading...'}</h1>
    </div>
    <div class="logo-controls">
      <label class="logo-upload-btn">
        {uploadingLogo ? 'Uploading...' : (logoUrl ? 'Change Logo' : 'Add Logo')}
        <input
          type="file"
          accept="image/png,image/jpeg,image/svg+xml"
          on:change={handleLogoUpload}
          disabled={uploadingLogo}
          style="display: none;"
        />
      </label>
      {#if logoUrl}
        <button class="logo-delete-btn" on:click={deleteLogo}>Remove Logo</button>
      {/if}
    </div>
  </div>

  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="content">
      <section class="section">
        <div class="section-header">
          <h2>Players in Group ({group.players?.length || 0})</h2>
          <div class="header-actions">
            <button class="toggle-weights-btn" on:click={() => showWeightBadges = !showWeightBadges}>
              {showWeightBadges ? 'Hide' : 'Show'} Weights
            </button>
            <button on:click={() => showAddPlayerModal = true}>+ Add Player</button>
          </div>
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
                  <span class="player-name">{player.name}</span>
                  {#if showWeightBadges}
                    <span class="weight-badge">Level {player.skill_weight} - {skillLevels[player.skill_weight].label}</span>
                  {/if}
                </div>
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
            <button
              class="separate-btn"
              on:click={separateSelectedPlayers}
              disabled={selectedPlayers.size === 0}
            >
              Separate Selected Players
            </button>
            <button
              class="remove-selected-btn"
              on:click={removeSelectedPlayers}
              disabled={selectedPlayers.size === 0}
            >
              Remove Selected Players
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

          {#if separatedGroups.length > 0}
            <div class="separated-groups">
              <h3>Separated Groups:</h3>
              {#each separatedGroups as separatedGroup, index}
                <div class="separated-group">
                  <span>{separatedGroup.map(id => getPlayerName(id)).join(', ')}</span>
                  <button class="btn-small-remove" on:click={() => removeSeparatedGroup(index)}>×</button>
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
          {#if numTeams === 2}
            <div class="jersey-option">
              <label class="checkbox-label">
                <input type="checkbox" bind:checked={useJerseyColors} />
                <span>Use jersey colors (Light/Dark)</span>
              </label>
            </div>
          {/if}
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
    </div>
  {/if}
</div>

<!-- Add Player Modal -->
{#if showAddPlayerModal}
  <div class="modal-overlay" on:click={() => showAddPlayerModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Add Players to Group</h2>
      {#if availablePlayers.length === 0}
        <p class="empty">No available players. Create players from the dashboard first!</p>
        <button on:click={() => showAddPlayerModal = false}>Close</button>
      {:else}
        <div class="player-list">
          {#each availablePlayers as player}
            <div class="player-list-item selectable" class:selected={selectedPlayersToAdd.has(player.id)} on:click={() => togglePlayerToAdd(player.id)}>
              <div class="checkbox">
                {#if selectedPlayersToAdd.has(player.id)}
                  <span class="checkmark">✓</span>
                {/if}
              </div>
              <div class="player-list-info">
                <span>{player.name}</span>
                <span class="weight-badge">Level {player.skill_weight}</span>
              </div>
            </div>
          {/each}
        </div>
        <div class="modal-actions">
          <button class="btn-primary" on:click={addPlayersToGroup} disabled={selectedPlayersToAdd.size === 0}>
            Add {selectedPlayersToAdd.size > 0 ? `${selectedPlayersToAdd.size} Player${selectedPlayersToAdd.size > 1 ? 's' : ''}` : 'Players'}
          </button>
          <button type="button" on:click={() => showAddPlayerModal = false}>Cancel</button>
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
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .back-btn {
    background: none;
    border: none;
    color: #2196F3;
    font-size: 16px;
    cursor: pointer;
    padding: 5px 0;
    align-self: flex-start;
  }

  .back-btn:hover {
    text-decoration: underline;
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

  .logo-controls {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  .logo-upload-btn {
    padding: 8px 16px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .logo-upload-btn:hover {
    background-color: #45a049;
  }

  .logo-delete-btn {
    padding: 8px 16px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .logo-delete-btn:hover {
    background-color: #da190b;
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

  .header-actions {
    display: flex;
    gap: 10px;
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

  .toggle-weights-btn {
    background-color: #9E9E9E !important;
  }

  .toggle-weights-btn:hover {
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
    padding: 8px 12px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    transition: all 0.2s;
    margin-bottom: 6px;
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
    gap: 12px;
    flex: 1;
    cursor: pointer;
  }

  .checkbox {
    width: 20px;
    height: 20px;
    border: 2px solid #bdbdbd;
    border-radius: 3px;
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
    font-size: 14px;
  }

  .player-name {
    font-size: 15px;
    font-weight: 500;
    color: #333;
  }

  .weight-badge {
    display: inline-block;
    background-color: #e3f2fd;
    color: #1976d2;
    padding: 3px 8px;
    border-radius: 3px;
    font-size: 12px;
    font-weight: 500;
  }

  .btn-remove {
    padding: 5px 10px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 3px;
    cursor: pointer;
    font-size: 13px;
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

  .separate-btn {
    padding: 10px 20px;
    background-color: #9C27B0;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .separate-btn:hover:not(:disabled) {
    background-color: #7B1FA2;
  }

  .separate-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .remove-selected-btn {
    padding: 10px 20px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .remove-selected-btn:hover:not(:disabled) {
    background-color: #da190b;
  }

  .remove-selected-btn:disabled {
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

  .separated-groups {
    margin-top: 20px;
    padding: 15px;
    background-color: #f3e5f5;
    border-radius: 4px;
  }

  .separated-groups h3 {
    margin: 0 0 10px 0;
    font-size: 16px;
    color: #6a1b9a;
  }

  .separated-group {
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

  .jersey-option {
    margin-top: 20px;
    padding: 15px;
    background-color: #f9f9f9;
    border-radius: 4px;
    border: 1px solid #e0e0e0;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    font-size: 16px;
    color: #333;
  }

  .checkbox-label input[type="checkbox"] {
    width: 20px;
    height: 20px;
    cursor: pointer;
  }

  .checkbox-label span {
    user-select: none;
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
    align-items: center;
    gap: 15px;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    margin-bottom: 10px;
    transition: all 0.2s;
  }

  .player-list-item.selectable {
    cursor: pointer;
  }

  .player-list-item.selectable:hover {
    background-color: #f5f5f5;
    border-color: #4CAF50;
  }

  .player-list-item.selected {
    border-color: #4CAF50;
    background-color: #f1f8f4;
  }

  .player-list-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex: 1;
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

  .btn-primary:hover:not(:disabled) {
    background-color: #45a049;
  }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .modal-actions button[type="button"] {
    background-color: #e0e0e0;
    color: #333;
  }

  .modal-actions button[type="button"]:hover {
    background-color: #d0d0d0;
  }

  @media (max-width: 768px) {
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

    .header-actions {
      width: 100%;
      flex-direction: column;
    }

    .section-header button {
      width: 100%;
    }
  }
</style>
