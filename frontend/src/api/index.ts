import { LeaderboardRow, PlayerProfile, GameResult, PlayerCreate } from '../types';

const API_BASE_URL = 'https://api.luincpong.com';

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new ApiError(response.status, errorText || response.statusText);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    throw error;
  }
}

export const api = {
  // GET /leaderboard (not /leaderboards)
  getLeaderboard: (): Promise<LeaderboardRow[]> =>
    fetchApi<LeaderboardRow[]>('/leaderboard'),

  // GET /players/:id
  getPlayer: (id: number): Promise<PlayerProfile> =>
    fetchApi<PlayerProfile>(`/players/${id}`),

  // POST /players
  addPlayer: (data: PlayerCreate): Promise<void> =>
    fetchApi<void>('/players', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // POST /games
  addGame: (data: GameResult): Promise<void> =>
    fetchApi<void>('/games', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
};
