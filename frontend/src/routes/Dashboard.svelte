<script>
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import { authAPI, playersAPI, groupsAPI } from '../lib/api';
  import { skillLevels } from '../lib/store';

  let players = [];
  let groups = [];
  let loading = true;
  let showNewPlayerModal = false;
  let showNewGroupModal = false;
  let showEditPlayerModal = false;
  let editingPlayer = null;

  let newPlayerName = '';
  let newPlayerWeight = 3;
  let newGroupName = '';

  onMount(async () => {
    if (!localStorage.getItem('token')) {
      navigate('/login');
      return;
    }

    try {
      await authAPI.me();
      await loadData();
    } catch (err) {
      navigate('/login');
    }
  });

  async function loadData() {
    loading = true;
    try {
      [players, groups] = await Promise.all([
        playersAPI.getAll(),
        groupsAPI.getAll()
      ]);
    } catch (err) {
      console.error('Failed to load data:', err);
    } finally {
      loading = false;
    }
  }

  async function createPlayer() {
    try {
      await playersAPI.create(newPlayerName, newPlayerWeight);
      newPlayerName = '';
      newPlayerWeight = 3;
      showNewPlayerModal = false;
      await loadData();
    } catch (err) {
      alert(err.message);
    }
  }

  async function createGroup() {
    try {
      await groupsAPI.create(newGroupName);
      newGroupName = '';
      showNewGroupModal = false;
      await loadData();
    } catch (err) {
      alert(err.message);
    }
  }

  async function deletePlayer(id) {
    if (confirm('Are you sure you want to delete this player?')) {
      try {
        await playersAPI.delete(id);
        await loadData();
      } catch (err) {
        alert(err.message);
      }
    }
  }

  async function deleteGroup(id) {
    if (confirm('Are you sure you want to delete this group?')) {
      try {
        await groupsAPI.delete(id);
        await loadData();
      } catch (err) {
        alert(err.message);
      }
    }
  }

  function openEditPlayer(player) {
    editingPlayer = { ...player };
    showEditPlayerModal = true;
  }

  async function updatePlayer() {
    try {
      await playersAPI.update(editingPlayer.id, editingPlayer.name, editingPlayer.skill_weight);
      showEditPlayerModal = false;
      editingPlayer = null;
      await loadData();
    } catch (err) {
      alert(err.message);
    }
  }

  function logout() {
    authAPI.logout();
  }
</script>

<div class="container">
  <header>
    <h1>Stick Toss</h1>
    <button class="logout-btn" on:click={logout}>Logout</button>
  </header>

  {#if loading}
    <div class="loading">Loading...</div>
  {:else}
    <div class="content">
      <section class="section">
        <div class="section-header">
          <h2>Players</h2>
          <button on:click={() => showNewPlayerModal = true}>+ Add Player</button>
        </div>

        {#if players.length === 0}
          <p class="empty">No players yet. Create your first player!</p>
        {:else}
          <div class="grid">
            {#each players as player}
              <div class="card">
                <div class="card-content">
                  <h3>{player.name}</h3>
                  <div class="weight-badge">
                    Level {player.skill_weight} - {skillLevels[player.skill_weight].label}
                  </div>
                </div>
                <div class="card-actions">
                  <button class="btn-small" on:click={() => openEditPlayer(player)}>Edit</button>
                  <button class="btn-small btn-danger" on:click={() => deletePlayer(player.id)}>Delete</button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </section>

      <section class="section">
        <div class="section-header">
          <h2>Groups</h2>
          <button on:click={() => showNewGroupModal = true}>+ Add Group</button>
        </div>

        {#if groups.length === 0}
          <p class="empty">No groups yet. Create your first group!</p>
        {:else}
          <div class="grid">
            {#each groups as group}
              <div class="card">
                <div class="card-content">
                  <h3>{group.name}</h3>
                </div>
                <div class="card-actions">
                  <button class="btn-small" on:click={() => navigate(`/group/${group.id}`)}>Open</button>
                  <button class="btn-small btn-danger" on:click={() => deleteGroup(group.id)}>Delete</button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </section>
    </div>
  {/if}
</div>

<!-- New Player Modal -->
{#if showNewPlayerModal}
  <div class="modal-overlay" on:click={() => showNewPlayerModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Add New Player</h2>
      <form on:submit|preventDefault={createPlayer}>
        <div class="form-group">
          <label>Name</label>
          <input type="text" bind:value={newPlayerName} required />
        </div>
        <div class="form-group">
          <label>
            Skill Level
            <span class="info-icon" title={skillLevels[newPlayerWeight].description}>ℹ️</span>
          </label>
          <select bind:value={newPlayerWeight}>
            {#each Object.entries(skillLevels) as [level, info]}
              <option value={parseInt(level)}>{level} - {info.label}</option>
            {/each}
          </select>
          <p class="skill-description">{skillLevels[newPlayerWeight].description}</p>
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showNewPlayerModal = false}>Cancel</button>
          <button type="submit">Create</button>
        </div>
      </form>
    </div>
  </div>
{/if}

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
          <select bind:value={editingPlayer.skill_weight}>
            {#each Object.entries(skillLevels) as [level, info]}
              <option value={parseInt(level)}>{level} - {info.label}</option>
            {/each}
          </select>
          <p class="skill-description">{skillLevels[editingPlayer.skill_weight].description}</p>
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showEditPlayerModal = false}>Cancel</button>
          <button type="submit">Update</button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- New Group Modal -->
{#if showNewGroupModal}
  <div class="modal-overlay" on:click={() => showNewGroupModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Add New Group</h2>
      <form on:submit|preventDefault={createGroup}>
        <div class="form-group">
          <label>Group Name</label>
          <input type="text" bind:value={newGroupName} placeholder="Tuesday Night Hockey" required />
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showNewGroupModal = false}>Cancel</button>
          <button type="submit">Create</button>
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
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 40px;
  }

  h1 {
    font-size: 36px;
    color: #333;
    margin: 0;
  }

  .logout-btn {
    padding: 8px 16px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .logout-btn:hover {
    background-color: #da190b;
  }

  .loading {
    text-align: center;
    padding: 40px;
    font-size: 18px;
    color: #666;
  }

  .section {
    margin-bottom: 40px;
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
    font-size: 14px;
  }

  .section-header button:hover {
    background-color: #45a049;
  }

  .empty {
    text-align: center;
    padding: 40px;
    color: #999;
    font-style: italic;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 20px;
  }

  .card-content h3 {
    margin: 0 0 10px 0;
    font-size: 18px;
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

  .card-actions {
    margin-top: 15px;
    display: flex;
    gap: 10px;
  }

  .btn-small {
    padding: 6px 12px;
    background-color: #2196F3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .btn-small:hover {
    background-color: #0b7dda;
  }

  .btn-danger {
    background-color: #f44336;
  }

  .btn-danger:hover {
    background-color: #da190b;
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
    margin-top: 10px;
    padding: 10px;
    background-color: #f5f5f5;
    border-radius: 4px;
    font-size: 13px;
    color: #666;
    line-height: 1.4;
  }

  .modal-actions {
    display: flex;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .modal-actions button {
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .modal-actions button[type="button"] {
    background-color: #e0e0e0;
    color: #333;
  }

  .modal-actions button[type="submit"] {
    background-color: #4CAF50;
    color: white;
  }

  .modal-actions button:hover {
    opacity: 0.9;
  }

  @media (max-width: 768px) {
    .grid {
      grid-template-columns: 1fr;
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
