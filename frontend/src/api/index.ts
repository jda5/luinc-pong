import { PlayerProfile, GameResult, PlayerCreate, Achievement, IndexPageData, HeadToHead } from '../types';

// const API_BASE_URL = 'https://api.luincpong.com';
const API_BASE_URL = 'http://localhost:8080';

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
  // GET / - Index page data with leaderboard and global stats
  getIndexPageData: (): Promise<IndexPageData> =>
    fetchApi<IndexPageData>('/'),

  // GET /achievements
  getAchievements: (): Promise<Achievement[]> =>
    fetchApi<Achievement[]>('/achievements'),

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

  // GET /head-to-head?p1=:p1&p2=:p2
  getHeadToHead: (p1: number, p2: number): Promise<HeadToHead> =>
    fetchApi<HeadToHead>(`/head-to-head?p1=${p1}&p2=${p2}`),
};
