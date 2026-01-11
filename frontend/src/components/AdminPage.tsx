import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Trash2, Shield, ChevronLeft } from 'lucide-react';
import { api } from '../api';
import type { Game } from '../types';

// Utility function for date formatting
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-GB', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).split('/').reverse().join('-');
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

// Password Gate Component
interface PasswordGateProps {
  onAuthenticated: () => void;
}

const PasswordGate: React.FC<PasswordGateProps> = ({ onAuthenticated }) => {
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isShaking, setIsShaking] = useState(false);

  // Simple client-side password - in production, use proper auth
  const ADMIN_PASSWORD = 'iamanumpty';

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password === ADMIN_PASSWORD) {
      sessionStorage.setItem('adminAuthenticated', 'true');
      onAuthenticated();
    } else {
      setError('Incorrect password');
      setIsShaking(true);
      setTimeout(() => setIsShaking(false), 500);
    }
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-4">
      <div className={`athletic-container max-w-md w-full ${isShaking ? 'animate-shake' : ''}`}>
        <div className="text-center mb-8">
          <h1 className="athletic-heading text-2xl mb-2">Admin Portal</h1>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="athletic-label block mb-2">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                setError('');
              }}
              className="athletic-input"
              placeholder="Enter password..."
              autoFocus
            />
          </div>

          {error && (
            <div className="text-red-600 text-sm text-center font-medium">
              {error}
            </div>
          )}

          <div className="flex gap-4 pt-2">
            <Link
              to="/"
              className="athletic-btn athletic-btn-secondary flex-1 flex items-center justify-center gap-2"
            >
              <ChevronLeft className="w-4 h-4" />
              Home
            </Link>
            <button
              type="submit"
              className="athletic-btn athletic-btn-primary flex-1 flex items-center justify-center gap-2"
            >
              <Shield className="w-4 h-4" />
              Enter
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

// Delete Confirmation Modal
interface DeleteModalProps {
  game: Game;
  onConfirm: () => void;
  onCancel: () => void;
  isDeleting: boolean;
}

const DeleteModal: React.FC<DeleteModalProps> = ({ game, onConfirm, onCancel, isDeleting }) => {
  const formatScore = () => {
    if (game.winnerScore !== null && game.loserScore !== null) {
      return `${game.winnerScore} - ${game.loserScore}`;
    }
    return 'No score recorded';
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div className="athletic-modal max-w-lg">
        <div className="athletic-modal-header">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
              <Trash2 className="w-5 h-5 text-red-600" />
            </div>
            <h2 className="athletic-modal-title">Delete Game</h2>
          </div>
        </div>

        <div className="athletic-modal-body">
          <p className="athletic-body text-grey-900 mb-6">
            Are you sure you want to delete this game? This action cannot be undone.
          </p>

          <div className="admin-game-card-preview mb-6">
            <div className="flex items-center justify-between mb-3">
              <span className="athletic-heading text-lg">{game.winner.name}</span>
              <span className="text-lime-green font-bold">VS</span>
              <span className="athletic-heading text-lg">{game.loser.name}</span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <span className="athletic-label">Score</span>
              <span className="athletic-body font-mono">{formatScore()}</span>
            </div>
            <div className="flex items-center justify-between text-sm mt-2">
              <span className="athletic-label">Date</span>
              <span className="athletic-body">{formatDate(game.createdAt)}</span>
            </div>
          </div>

          <div className="flex gap-3">
            <button
              onClick={onCancel}
              className="athletic-btn athletic-btn-secondary flex-1"
              disabled={isDeleting}
            >
              Cancel
            </button>
            <button
              onClick={onConfirm}
              className="athletic-btn athletic-btn-danger flex-1"
              disabled={isDeleting}
            >
              {isDeleting ? (
                <div className="athletic-spinner w-4 h-4"></div>
              ) : (
                <>
                  <Trash2 className="w-4 h-4" />
                  Delete Game
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

// Game Card Component
interface GameCardProps {
  game: Game;
  onDelete: (game: Game) => void;
}

const GameCard: React.FC<GameCardProps> = ({ game, onDelete }) => {
  const formatScore = () => {
    if (game.winnerScore !== null && game.loserScore !== null) {
      return `${game.winnerScore} - ${game.loserScore}`;
    }
    return null;
  };

  const score = formatScore();

  return (
    <div className="admin-game-card">
      <div className="admin-game-card-header">
        <div className="admin-game-players">
          <div className="admin-game-player winner">
            <span className="admin-game-badge">W</span>
            <span className="admin-game-player-name">{game.winner.name}</span>
          </div>
          <div className="admin-game-vs">VS</div>
          <div className="admin-game-player loser">
            <span className="admin-game-player-name">{game.loser.name}</span>
            <span className="admin-game-badge loser">L</span>
          </div>
        </div>
        <button
          onClick={() => onDelete(game)}
          className="admin-delete-btn"
          title="Delete game"
        >
          <Trash2 className="w-4 h-4" />
        </button>
      </div>
      
      <div className="admin-game-card-footer">
        {score && (
          <div className="admin-game-score">
            <span className="athletic-label">Score</span>
            <span className="admin-game-score-value">{score}</span>
          </div>
        )}
        <div className="admin-game-date">
          <span className="athletic-label">Played</span>
          <span className="admin-game-date-value">{formatRelativeTime(game.createdAt)}</span>
        </div>
      </div>
    </div>
  );
};

// Admin Dashboard Component
const AdminDashboard: React.FC = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [gameToDelete, setGameToDelete] = useState<Game | null>(null);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const queryClient = useQueryClient();

  const GAMES_PER_PAGE = 50;

  // Initial load query
  const { isLoading, error } = useQuery({
    queryKey: ['adminGames', 1],
    queryFn: () => api.getGames(1),
    staleTime: 0,
  });

  // Get the cached initial games data
  const initialGames = queryClient.getQueryData<Game[]>(['adminGames', 1]);

  // Update games when initial data loads
  React.useEffect(() => {
    if (initialGames && initialGames.length > 0 && games.length === 0) {
      setGames(initialGames);
      setHasMore(initialGames.length === GAMES_PER_PAGE);
    } else if (initialGames && initialGames.length === 0 && games.length === 0) {
      // Handle case where there are no games
      setHasMore(false);
    }
  }, [initialGames]);

  // Load more handler
  const handleLoadMore = async () => {
    const nextPage = page + 1;
    setIsLoadingMore(true);
    try {
      const newGames = await api.getGames(nextPage);
      setGames(prev => [...prev, ...newGames]);
      setPage(nextPage);
      setHasMore(newGames.length === GAMES_PER_PAGE);
    } catch (err) {
      console.error('Failed to load more games:', err);
    } finally {
      setIsLoadingMore(false);
    }
  };

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: (gameId: number) => api.deleteGame(gameId),
    onSuccess: (_, deletedGameId) => {
      setGames(prev => prev.filter(g => g.id !== deletedGameId));
      setGameToDelete(null);
      // Invalidate all related queries
      queryClient.invalidateQueries({ queryKey: ['indexPageData'] });
      queryClient.invalidateQueries({ queryKey: ['player'] });
      queryClient.invalidateQueries({ queryKey: ['adminGames'] });
    },
    onError: (error) => {
      console.error('Failed to delete game:', error);
      alert('Failed to delete game. Please try again.');
    },
  });

  const handleDeleteClick = (game: Game) => {
    setGameToDelete(game);
  };

  const handleConfirmDelete = () => {
    if (gameToDelete) {
      deleteMutation.mutate(gameToDelete.id);
    }
  };

  const handleCancelDelete = () => {
    setGameToDelete(null);
  };

  if (isLoading && games.length === 0) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="athletic-spinner mx-auto mb-4"></div>
          <p className="athletic-body text-grey-300">Loading games...</p>
        </div>
      </div>
    );
  }

  if (error && games.length === 0) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="athletic-container max-w-md text-center">
          <h2 className="athletic-heading text-2xl mb-4">Error Loading Games</h2>
          <p className="athletic-body text-grey-300 mb-4">
            Failed to load games. Please try again.
          </p>
          <button 
            onClick={() => window.location.reload()}
            className="athletic-btn athletic-btn-primary"
          >
            Retry
          </button>
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
            <div className="flex items-center gap-4">
              <Link to="/" className="athletic-btn athletic-btn-secondary p-2">
                <ChevronLeft className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="athletic-display text-3xl flex items-center gap-3">
                  Admin Dashboard
                </h1>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
     
        {/* Games Section */}
        <div className="athletic-container">
          <div className="mb-6">
            <h2 className="athletic-heading text-2xl mb-2">Game History</h2>
          </div>

          {games.length > 0 ? (
            <>
              <div className="admin-games-grid">
                {games.map((game) => (
                  <GameCard
                    key={game.id}
                    game={game}
                    onDelete={handleDeleteClick}
                  />
                ))}
              </div>

              {/* Load More Button */}
              {hasMore && (
                <div className="mt-12 text-center">
                  <button
                    onClick={handleLoadMore}
                    className="athletic-btn athletic-btn-primary"
                    disabled={isLoadingMore}
                  >
                    {isLoadingMore ? (
                      <>
                        <div className="athletic-spinner w-4 h-4"></div>
                        Loading...
                      </>
                    ) : (
                      'Load More Games'
                    )}
                  </button>
                </div>
              )}

              {!hasMore && games.length > 0 && (
                <div className="mt-8 text-center">
                  <p className="athletic-label">All games loaded</p>
                </div>
              )}
            </>
          ) : (
            <div className="text-center py-12">
              <p className="athletic-body text-grey-300 text-lg">
                No games found
              </p>
            </div>
          )}
        </div>
      </main>

      {/* Delete Confirmation Modal */}
      {gameToDelete && (
        <DeleteModal
          game={gameToDelete}
          onConfirm={handleConfirmDelete}
          onCancel={handleCancelDelete}
          isDeleting={deleteMutation.isPending}
        />
      )}
    </div>
  );
};

// Main Admin Page Component
const AdminPage: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return sessionStorage.getItem('adminAuthenticated') === 'true';
  });

  if (!isAuthenticated) {
    return <PasswordGate onAuthenticated={() => setIsAuthenticated(true)} />;
  }

  return <AdminDashboard />;
};

export default AdminPage;
