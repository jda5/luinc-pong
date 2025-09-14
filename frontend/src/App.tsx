import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useParams } from 'react-router-dom';
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from './api';
import type { GameResult, PlayerCreate } from './types';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      refetchOnWindowFocus: false,
    },
  },
});

// Utility functions
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).replace(/\//g, '/');
};

const formatRelativeTime = (dateString: string) => {
  const now = new Date();
  const date = new Date(dateString);
  const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));
  
  if (diffInHours < 1) return 'just now';
  if (diffInHours < 24) return `${diffInHours}h ago`;
  
  const diffInDays = Math.floor(diffInHours / 24);
  if (diffInDays < 7) return `${diffInDays}d ago`;
  
  return formatDate(dateString);
};

// Modal prop interfaces
interface AddPlayerModalProps {
  onClose: () => void;
}

interface AddGameModalProps {
  onClose: () => void;
}

// Add Player Modal
const AddPlayerModal: React.FC<AddPlayerModalProps> = ({ onClose }) => {
  const [name, setName] = useState('');
  const queryClient = useQueryClient();

  const addPlayerMutation = useMutation<void, Error, PlayerCreate>({
    mutationFn: (data: PlayerCreate) => api.addPlayer(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['leaderboard'] });
      onClose();
      setName('');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim()) {
      addPlayerMutation.mutate({
        name: name.trim(),
      });
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="athletic-modal">
        <div className="athletic-modal-header">
          <div className="flex items-center justify-between">
            <h2 className="athletic-modal-title">Add New Player</h2>
            <button
              onClick={onClose}
              className="athletic-btn athletic-btn-secondary text-sm px-3 py-2"
            >
              ✕
            </button>
          </div>
        </div>

        <div className="athletic-modal-body">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="athletic-label block mb-2">Player Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="athletic-input"
                placeholder="Enter player name"
                required
              />
            </div>

            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={onClose}
                className="athletic-btn athletic-btn-secondary flex-1"
                disabled={addPlayerMutation.isPending}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="athletic-btn athletic-btn-primary flex-1"
                disabled={addPlayerMutation.isPending || !name.trim()}
              >
                {addPlayerMutation.isPending ? (
                  <div className="athletic-spinner w-4 h-4"></div>
                ) : (
                  'Add Player'
                )}
              </button>
            </div>

            {addPlayerMutation.error && (
              <div className="text-red-600 text-sm text-center">
                Failed to add player. Please try again.
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
};

// Add Game Modal
const AddGameModal: React.FC<AddGameModalProps> = ({ onClose }) => {
  const [winnerId, setWinnerId] = useState('');
  const [loserId, setLoserId] = useState('');
  const [winnerScore, setWinnerScore] = useState('');
  const [loserScore, setLoserScore] = useState('');
  const queryClient = useQueryClient();

  const { data: leaderboard } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: api.getLeaderboard,
  });

  const addGameMutation = useMutation<void, Error, GameResult>({
    mutationFn: (data: GameResult) => api.addGame(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['leaderboard'] });
      queryClient.invalidateQueries({ queryKey: ['player'] });
      onClose();
      setWinnerId('');
      setLoserId('');
      setWinnerScore('');
      setLoserScore('');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (winnerId && loserId && winnerId !== loserId) {
      const gameData: GameResult = {
        winnerId: parseInt(winnerId),
        loserId: parseInt(loserId),
      };

      // Add scores if provided
      if (winnerScore.trim()) {
        gameData.winnerScore = parseInt(winnerScore);
      }
      if (loserScore.trim()) {
        gameData.loserScore = parseInt(loserScore);
      }

      addGameMutation.mutate(gameData);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="athletic-modal">
        <div className="athletic-modal-header">
          <div className="flex items-center justify-between">
            <h2 className="athletic-modal-title">Record Game</h2>
            <button
              onClick={onClose}
              className="athletic-btn athletic-btn-secondary text-sm px-3 py-2"
            >
              ✕
            </button>
          </div>
        </div>

        <div className="athletic-modal-body">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="athletic-label block mb-2">Winner</label>
              <select
                value={winnerId}
                onChange={(e) => setWinnerId(e.target.value)}
                className="athletic-select"
                required
              >
                <option value="">Select winner</option>
                {leaderboard?.map((player) => (
                  <option key={player.id} value={player.id}>
                    {player.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="athletic-label block mb-2">Loser</label>
              <select
                value={loserId}
                onChange={(e) => setLoserId(e.target.value)}
                className="athletic-select"
                required
              >
                <option value="">Select loser</option>
                {leaderboard?.map((player) => (
                  <option key={player.id} value={player.id}>
                    {player.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="athletic-label block mb-2">Winner Score (Optional)</label>
                <input
                  type="number"
                  value={winnerScore}
                  onChange={(e) => setWinnerScore(e.target.value)}
                  className="athletic-input"
                  placeholder="e.g. 11"
                  min="0"
                  max="255"
                />
              </div>
              <div>
                <label className="athletic-label block mb-2">Loser Score (Optional)</label>
                <input
                  type="number"
                  value={loserScore}
                  onChange={(e) => setLoserScore(e.target.value)}
                  className="athletic-input"
                  placeholder="e.g. 9"
                  min="0"
                  max="255"
                />
              </div>
            </div>

            <div className="flex gap-3 pt-4">
              <button
                type="button"
                onClick={onClose}
                className="athletic-btn athletic-btn-secondary flex-1"
                disabled={addGameMutation.isPending}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="athletic-btn athletic-btn-primary flex-1"
                disabled={addGameMutation.isPending || !winnerId || !loserId || winnerId === loserId}
              >
                {addGameMutation.isPending ? (
                  <div className="athletic-spinner w-4 h-4"></div>
                ) : (
                  'Record Game'
                )}
              </button>
            </div>

            {winnerId === loserId && winnerId && (
              <div className="text-red-600 text-sm text-center">
                Winner and loser must be different players.
              </div>
            )}

            {addGameMutation.error && (
              <div className="text-red-600 text-sm text-center">
                {addGameMutation.error.message || 'Failed to record game. Please try again.'}
              </div>
            )}
          </form>
        </div>
      </div>
    </div>
  );
};


// Home/Leaderboard Page
const HomePage: React.FC = () => {
  const [showAddPlayer, setShowAddPlayer] = useState(false);
  const [showAddGame, setShowAddGame] = useState(false);

  const { data: leaderboard, isLoading, error } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: api.getLeaderboard,
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="athletic-spinner mx-auto mb-4"></div>
          <p className="athletic-body text-grey-300">Loading leaderboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    const isNetworkError = errorMessage.includes('fetch') || errorMessage.includes('NetworkError') || errorMessage.includes('Failed to fetch');
    
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="athletic-container max-w-md text-center">
          <h2 className="athletic-heading text-2xl mb-4">
            {isNetworkError ? 'Backend Connection Error' : 'API Error'}
          </h2>
          <p className="athletic-body text-grey-300 mb-4">
            {isNetworkError 
              ? 'Cannot connect to the backend server at http://localhost:8080' 
              : 'There was an error loading the leaderboard'}
          </p>
          <p className="athletic-label mb-6 text-red-600">
            {errorMessage}
          </p>
          <div className="flex gap-3">
            <button 
              onClick={() => window.location.reload()}
              className="athletic-btn athletic-btn-primary flex-1"
            >
              Retry
            </button>
            <button
              onClick={() => setShowAddPlayer(true)}
              className="athletic-btn athletic-btn-secondary flex-1"
            >
              Add Player Anyway
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="border-b border-grey-100 bg-white sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between">
            <h1 className="athletic-display text-4xl">LUINC PONG</h1>
            <div className="flex gap-3">
              <button
                onClick={() => setShowAddPlayer(true)}
                className="athletic-btn athletic-btn-secondary"
              >
                Add Player
              </button>
              <button
                onClick={() => setShowAddGame(true)}
                className="athletic-btn athletic-btn-primary"
              >
                Record Game
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="athletic-container">
          <div className="mb-8">
            <h2 className="athletic-heading text-2xl mb-2">Leaderboard</h2>
            <p className="athletic-label">Current Rankings</p>
          </div>

          {/* Leaderboard Grid */}
          <div className="space-y-0">
            {/* Header Row */}
            <div className="leaderboard-row border-b-2 border-grey-200 bg-grey-50">
              <div className="athletic-label">RANK</div>
              <div className="athletic-label">PLAYER</div>
              <div className="athletic-label text-right">ELO</div>
            </div>

            {/* Player Rows */}
            {leaderboard && leaderboard.length > 0 ? (
              leaderboard.map((player, index) => (
                <Link
                  key={player.id}
                  to={`/player/${player.id}`}
                  className="leaderboard-row block hover:no-underline"
                >
                  <div className="leaderboard-rank">
                    {index + 1}
                  </div>
                  <div className="leaderboard-name">
                    {player.name || 'Unknown Player'}
                  </div>
                  <div className="leaderboard-elo">
                    {player.eloRating != null ? Math.round(player.eloRating) : '1000'}
                  </div>
                </Link>
              ))
            ) : (
              <div className="py-12 text-center">
                <div className="athletic-container max-w-md mx-auto">
                  <h3 className="athletic-heading text-xl mb-4">No Players Yet</h3>
                  <p className="athletic-body text-grey-300 mb-6">
                    Add your first player to get started with the leaderboard.
                  </p>
                  <button
                    onClick={() => setShowAddPlayer(true)}
                    className="athletic-btn athletic-btn-primary"
                  >
                    Add First Player
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>

      {/* Modals */}
      {showAddPlayer && (
        <AddPlayerModal onClose={() => setShowAddPlayer(false)} />
      )}
      {showAddGame && (
        <AddGameModal onClose={() => setShowAddGame(false)} />
      )}
    </div>
  );
};

// Player Profile Page
const PlayerPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const playerId = parseInt(id || '0');

  const { data: player, isLoading, error } = useQuery({
    queryKey: ['player', playerId],
    queryFn: () => api.getPlayer(playerId),
    enabled: !!playerId,
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="athletic-spinner mx-auto mb-4"></div>
          <p className="athletic-body text-grey-300">Loading player profile...</p>
        </div>
      </div>
    );
  }

  if (error || !player) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="athletic-container max-w-md text-center">
          <h2 className="athletic-heading text-2xl mb-4">Player Not Found</h2>
          <p className="athletic-body text-grey-300 mb-6">
            The requested player could not be found.
          </p>
          <Link to="/" className="athletic-btn athletic-btn-primary">
            Back to Leaderboard
          </Link>
        </div>
      </div>
    );
  }

  const winPercentage = player.gamesPlayed > 0 
    ? Math.round((player.gamesWon / player.gamesPlayed) * 100) 
    : 0;

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="border-b border-grey-100 bg-white">
        <div className="max-w-6xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between">
            <Link to="/" className="athletic-btn athletic-btn-secondary">
              ← Back
            </Link>
            <h1 className="athletic-display text-4xl">LUINC PONG</h1>
            <div></div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
        {/* Player Header */}
        <div className="mb-8">
          <h1 className="athletic-display text-5xl mb-2">
            {player.name}
          </h1>
          <p className="athletic-label">PLAYER PROFILE</p>
        </div>

        {/* Personal Records Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <div className="pr-stat">
            <div className="pr-value">{player.eloRating != null ? Math.round(player.eloRating) : '1000'}</div>
            <div className="pr-label">ELO Rating</div>
          </div>
          <div className="pr-stat">
            <div className="pr-value">{winPercentage}%</div>
            <div className="pr-label">Win Rate</div>
          </div>
          <div className="pr-stat">
            <div className="pr-value">{player.gamesPlayed || 0}</div>
            <div className="pr-label">Games Played</div>
          </div>
        </div>

        {/* Recent Games */}
        <div className="athletic-container">
          <div className="mb-8">
            <h2 className="athletic-heading text-2xl mb-2">Recent Games</h2>
          </div>

          {player.recentGames && player.recentGames.length > 0 ? (
            <div className="space-y-0">
              {/* Header Row */}
              <div className="recent-games-item border-b-2 border-grey-200 bg-grey-50" style={{gridTemplateColumns: '80px 1fr 120px 100px'}}>
                <div className="athletic-label">RESULT</div>
                <div className="athletic-label">OPPONENT</div>
                <div className="athletic-label text-center">SCORE</div>
                <div className="athletic-label text-right">DATE</div>
              </div>

              {/* Game Rows */}
              {player.recentGames.map((game) => {
                const isWin = game.winner.id === player.id;
                const opponent = isWin ? game.loser : game.winner;
                
                // Format score display
                const formatScore = () => {
                  if (game.winnerScore !== null && game.loserScore !== null) {
                    if (isWin) {
                      return `${game.winnerScore} - ${game.loserScore}`;
                    } else {
                      return `${game.loserScore} - ${game.winnerScore}`;
                    }
                  }
                  return '';
                };
                
                return (
                  <div key={game.id} className="recent-games-item" style={{gridTemplateColumns: '80px 1fr 120px 100px'}}>
                    <div className={`game-result ${isWin ? 'win' : 'loss'}`}>
                      {isWin ? 'WIN' : 'LOSS'}
                    </div>
                    <div className="athletic-body">{opponent.name}</div>
                    <div className="athletic-body text-center font-mono">
                      {formatScore()}
                    </div>
                    <div className="athletic-body text-grey-300 text-right text-sm">
                      {formatRelativeTime(game.createdAt)}
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="athletic-body text-grey-300 text-lg">
                No games recorded
              </p>
              <p className="athletic-label mt-2">
                Start playing to build your recent games list
              </p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

// Main App Component
const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <div className="min-h-screen">
          <main>
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/player/:id" element={<PlayerPage />} />
            </Routes>
          </main>
        </div>
    </Router>
    </QueryClientProvider>
  );
};

export default App;
