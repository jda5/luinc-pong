export interface LeaderboardRow {
  id: number;
  name: string;
  eloRating: number;
}

export interface GlobalStats {
  totalGames: number;
  totalPoints: number;
}

export interface IndexPageData {
  leaderboard: LeaderboardRow[];
  globalStats: GlobalStats;
}

export interface Player {
  id: number;
  name: string;
}

export interface Game {
  id: number;
  winner: Player;
  loser: Player;
  winnerScore: number | null;
  loserScore: number | null;
  createdAt: string;
}

export interface Achievement {
  id: number;
  title: string;
  description: string;
}

export interface PlayerProfile {
  id: number;
  name: string;
  eloRating: number;
  createdAt: string;
  gamesPlayed: number;
  gamesWon: number;
  recentGames: Game[];
  achievements: Achievement[];
}

// Request body types
export interface GameResult {
  winnerId: number;
  loserId: number;
  winnerScore?: number;
  loserScore?: number;
}

export interface PlayerCreate {
  name: string;
}

// Head-to-Head types
export interface HeadToHeadPlayerStats {
  id: number;
  name: string;
  gamesWon: number;
  winProbability: number;
  longestWinStreak: number;
  totalPoints: number;
  avgPointsPerGame: number;
}

export interface ScoreStats {
  avgScoreDifferential: number;
  biggestBlowout: GameResult;
  mostCompetitive: GameResult;
}

export interface HeadToHead {
  player1: HeadToHeadPlayerStats;
  player2: HeadToHeadPlayerStats;
  firstPlayedAt: string;
  totalGameCount: number;
  recentGames: Game[];
  scoreStats: ScoreStats;
}
