const API_BASE = '/api';

// Helper to get auth token from localStorage
function getAuthToken() {
  return localStorage.getItem('token');
}

// Helper to make authenticated requests
async function fetchWithAuth(url, options = {}) {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    // Token expired or invalid, clear it
    localStorage.removeItem('token');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  return response;
}

// Auth API
export const authAPI = {
  async signup(email, password) {
    const response = await fetch(`${API_BASE}/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Signup failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  },

  async login(email, password) {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    return data;
  },

  async me() {
    const response = await fetchWithAuth(`${API_BASE}/auth/me`);
    if (!response.ok) {
      throw new Error('Failed to get user info');
    }
    return response.json();
  },

  logout() {
    localStorage.removeItem('token');
    window.location.href = '/login';
  },
};

// Players API
export const playersAPI = {
  async getAll() {
    const response = await fetchWithAuth(`${API_BASE}/players`);
    if (!response.ok) throw new Error('Failed to fetch players');
    return response.json();
  },

  async get(id) {
    const response = await fetchWithAuth(`${API_BASE}/players/${id}`);
    if (!response.ok) throw new Error('Failed to fetch player');
    return response.json();
  },

  async create(name, skillWeight) {
    const response = await fetchWithAuth(`${API_BASE}/players`, {
      method: 'POST',
      body: JSON.stringify({ name, skill_weight: skillWeight }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create player');
    }
    return response.json();
  },

  async update(id, name, skillWeight) {
    const response = await fetchWithAuth(`${API_BASE}/players/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name, skill_weight: skillWeight }),
    });
    if (!response.ok) throw new Error('Failed to update player');
    return response.json();
  },

  async delete(id) {
    const response = await fetchWithAuth(`${API_BASE}/players/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete player');
    return response.json();
  },
};

// Groups API
export const groupsAPI = {
  async getAll() {
    const response = await fetchWithAuth(`${API_BASE}/groups`);
    if (!response.ok) throw new Error('Failed to fetch groups');
    return response.json();
  },

  async get(id) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${id}`);
    if (!response.ok) throw new Error('Failed to fetch group');
    return response.json();
  },

  async create(name) {
    const response = await fetchWithAuth(`${API_BASE}/groups`, {
      method: 'POST',
      body: JSON.stringify({ name }),
    });
    if (!response.ok) throw new Error('Failed to create group');
    return response.json();
  },

  async update(id, name) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    });
    if (!response.ok) throw new Error('Failed to update group');
    return response.json();
  },

  async delete(id) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete group');
    return response.json();
  },

  async addPlayer(groupId, playerId) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${groupId}/players`, {
      method: 'POST',
      body: JSON.stringify({ player_id: playerId }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to add player to group');
    }
    return response.json();
  },

  async removePlayer(groupId, playerId) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${groupId}/players/${playerId}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to remove player from group');
    return response.json();
  },

  async generateTeams(groupId, numTeams, lockedPlayers = [], separatedPlayers = [], useJerseyColors = false) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${groupId}/generate-teams`, {
      method: 'POST',
      body: JSON.stringify({
        num_teams: numTeams,
        locked_players: lockedPlayers,
        separated_players: separatedPlayers,
        use_jersey_colors: useJerseyColors
      }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to generate teams');
    }
    return response.json();
  },

  async uploadLogo(groupId, file) {
    const token = getAuthToken();
    const formData = new FormData();
    formData.append('logo', file);

    const response = await fetch(`${API_BASE}/groups/${groupId}/logo`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: formData,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to upload logo');
    }
    return response.json();
  },

  async deleteLogo(groupId) {
    const response = await fetchWithAuth(`${API_BASE}/groups/${groupId}/logo`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete logo');
    }
    return response.json();
  },
};

// Game API (public, no auth)
export const gameAPI = {
  async get(shareId) {
    const response = await fetch(`${API_BASE}/game/${shareId}`);
    if (!response.ok) {
      throw new Error('Game not found');
    }
    return response.json();
  },
};
