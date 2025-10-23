import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../api';
import type { HeadToHead, LeaderboardRow } from '../types';

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

// Player Selection Component
interface PlayerSelectorProps {
  players: LeaderboardRow[];
  selectedPlayer1: number | null;
  selectedPlayer2: number | null;
  onPlayer1Change: (id: number) => void;
  onPlayer2Change: (id: number) => void;
}

const PlayerSelector: React.FC<PlayerSelectorProps> = ({
  players,
  selectedPlayer1,
  selectedPlayer2,
  onPlayer1Change,
  onPlayer2Change,
}) => {
  return (
    <div className="athletic-container mb-8">
      <div className="mb-6">
        <h2 className="athletic-heading text-2xl mb-2">Select Players</h2>
        <p className="athletic-label">Choose two players to compare head-to-head statistics</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label className="athletic-label block mb-2">Player 1</label>
          <select
            value={selectedPlayer1 || ''}
            onChange={(e) => onPlayer1Change(parseInt(e.target.value))}
            className="athletic-select"
          >
            <option value="">Select first player</option>
            {players.map((player) => (
              <option key={player.id} value={player.id}>
                {player.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="athletic-label block mb-2">Player 2</label>
          <select
            value={selectedPlayer2 || ''}
            onChange={(e) => onPlayer2Change(parseInt(e.target.value))}
            className="athletic-select"
          >
            <option value="">Select second player</option>
            {players.map((player) => (
              <option key={player.id} value={player.id}>
                {player.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      {selectedPlayer1 && selectedPlayer2 && selectedPlayer1 === selectedPlayer2 && (
        <div className="mt-4 text-red-600 text-sm text-center">
          Please select two different players.
        </div>
      )}
    </div>
  );
};

// Head-to-Head Statistics Display
interface HeadToHeadDisplayProps {
  data: HeadToHead;
}

const HeadToHeadDisplay: React.FC<HeadToHeadDisplayProps> = ({ data }) => {
  const { player1, player2, firstPlayedAt, totalGameCount, recentGames, scoreStats } = data;

  return (
    <div className="space-y-12">
      {/* Head-to-Head Header */}
      <div className="text-center mb-12">
        <div className="head-to-head-vs">
          <div className="head-to-head-vs-text">VS</div>
        </div>
        <h1 className="athletic-display text-4xl mb-2">
          {player1.name} vs {player2.name}
        </h1>
        <p className="athletic-label">HEAD-TO-HEAD STATISTICS</p>
        <p className="athletic-body text-grey-300 mt-2">
          First played: {formatDate(firstPlayedAt)} • Total games: {totalGameCount}
        </p>
      </div>

      {/* Player Stats Comparison */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
        {/* Player 1 Stats */}
        <div className={`player-stats-card ${player1.gamesWon > player2.gamesWon ? 'winner' : ''}`}>
          <div className="text-center mb-6">
            <h3 className="athletic-heading text-2xl mb-2">{player1.name}</h3>
            <p className="athletic-label">PLAYER 1</p>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="pr-stat">
              <div className="pr-value">{player1.gamesWon}</div>
              <div className="pr-label">WINS</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{Math.round(player1.winProbability * 100)}%</div>
              <div className="pr-label">WIN PROBABILITY</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{player1.longestWinStreak}</div>
              <div className="pr-label">BEST STREAK</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{Math.round(player1.avgPointsPerGame)}</div>
              <div className="pr-label">AVG POINTS</div>
            </div>
          </div>
        </div>

        {/* Player 2 Stats */}
        <div className={`player-stats-card ${player2.gamesWon > player1.gamesWon ? 'winner' : ''}`}>
          <div className="text-center mb-6">
            <h3 className="athletic-heading text-2xl mb-2">{player2.name}</h3>
            <p className="athletic-label">PLAYER 2</p>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="pr-stat">
              <div className="pr-value">{player2.gamesWon}</div>
              <div className="pr-label">WINS</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{Math.round(player2.winProbability * 100)}%</div>
              <div className="pr-label">WIN PROBABILITY</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{player2.longestWinStreak}</div>
              <div className="pr-label">BEST STREAK</div>
            </div>
            <div className="pr-stat">
              <div className="pr-value">{Math.round(player2.avgPointsPerGame)}</div>
              <div className="pr-label">AVG POINTS</div>
            </div>
          </div>
        </div>
      </div>

      {/* Score Statistics */}
      <div className="athletic-container mb-12">
        <div className="mb-6">
          <h2 className="athletic-heading text-2xl mb-2">Score Statistics</h2>
          <p className="athletic-label">GAME SCORING ANALYSIS</p>
        </div>

        <div className="score-stats-grid">
          <div className="score-stat-item">
            <div className="score-stat-value">{Math.round(scoreStats.avgScoreDifferential)}</div>
            <div className="score-stat-label">AVG SCORE DIFF</div>
          </div>
          
          <div className="score-stat-item">
            <div className="score-stat-value">
              {scoreStats.biggestBlowout.winnerScore !== undefined && scoreStats.biggestBlowout.loserScore !== undefined
                ? `${scoreStats.biggestBlowout.winnerScore}-${scoreStats.biggestBlowout.loserScore}`
                : 'N/A'
              }
            </div>
            <div className="score-stat-label">BIGGEST BLOWOUT</div>
          </div>
          
          <div className="score-stat-item">
            <div className="score-stat-value">
              {scoreStats.mostCompetitive.winnerScore !== undefined && scoreStats.mostCompetitive.loserScore !== undefined
                ? `${scoreStats.mostCompetitive.winnerScore}-${scoreStats.mostCompetitive.loserScore}`
                : 'N/A'
              }
            </div>
            <div className="score-stat-label">MOST COMPETITIVE</div>
          </div>
        </div>
      </div>

      {/* Recent Games */}
      <div className="athletic-container">
        <div className="mb-6">
          <h2 className="athletic-heading text-2xl mb-2">Recent Games</h2>
          <p className="athletic-label">LATEST MATCHUPS</p>
        </div>

        {recentGames.length > 0 ? (
          <div className="space-y-0">
            {/* Header Row */}
            <div className="recent-games-item border-b-2 border-grey-200 bg-grey-50" style={{gridTemplateColumns: '80px 1fr 120px 100px'}}>
              <div className="athletic-label">RESULT</div>
              <div className="athletic-label">OPPONENT</div>
              <div className="athletic-label text-center">SCORE</div>
              <div className="athletic-label text-right">DATE</div>
            </div>

            {/* Game Rows */}
            {recentGames.map((game) => {
              const isPlayer1Win = game.winner.id === player1.id;
              const opponent = isPlayer1Win ? game.loser : game.winner;
              
              // Format score display
              const formatScore = () => {
                if (game.winnerScore !== null && game.loserScore !== null) {
                  if (isPlayer1Win) {
                    return `${game.winnerScore} - ${game.loserScore}`;
                  } else {
                    return `${game.loserScore} - ${game.winnerScore}`;
                  }
                }
                return '';
              };
              
              return (
                <div key={game.id} className="recent-games-item" style={{gridTemplateColumns: '80px 1fr 120px 100px'}}>
                  <div className={`game-result ${isPlayer1Win ? 'win' : 'loss'}`}>
                    {isPlayer1Win ? 'WIN' : 'LOSS'}
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
              No games recorded between these players
            </p>
            <p className="athletic-label mt-2">
              Start playing to build the head-to-head history
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

// Main Head-to-Head Page Component
const HeadToHeadPage: React.FC = () => {
  const [selectedPlayer1, setSelectedPlayer1] = useState<number | null>(null);
  const [selectedPlayer2, setSelectedPlayer2] = useState<number | null>(null);

  // Get players list for selection
  const { data: indexData, isLoading: playersLoading, error: playersError } = useQuery({
    queryKey: ['indexPageData'],
    queryFn: api.getIndexPageData,
  });

  // Get head-to-head data when both players are selected
  const { data: headToHeadData, isLoading: headToHeadLoading, error: headToHeadError } = useQuery({
    queryKey: ['headToHead', selectedPlayer1, selectedPlayer2],
    queryFn: () => api.getHeadToHead(selectedPlayer1!, selectedPlayer2!),
    enabled: !!(selectedPlayer1 && selectedPlayer2 && selectedPlayer1 !== selectedPlayer2),
  });

  if (playersLoading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="athletic-spinner mx-auto mb-4"></div>
          <p className="athletic-body text-grey-300">Loading players...</p>
        </div>
      </div>
    );
  }

  if (playersError) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="athletic-container max-w-md text-center">
          <h2 className="athletic-heading text-2xl mb-4">Error Loading Players</h2>
          <p className="athletic-body text-grey-300 mb-6">
            There was an error loading the player list.
          </p>
          <Link to="/" className="athletic-btn athletic-btn-primary">
            Back to Leaderboard
          </Link>
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
            <Link to="/" className="athletic-btn athletic-btn-secondary">
              ← Back
            </Link>
            <h1 className="athletic-display text-4xl">⚔️ Head-to-Head</h1>
            <div></div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
        {/* Player Selection */}
        <PlayerSelector
          players={indexData?.leaderboard || []}
          selectedPlayer1={selectedPlayer1}
          selectedPlayer2={selectedPlayer2}
          onPlayer1Change={setSelectedPlayer1}
          onPlayer2Change={setSelectedPlayer2}
        />

        {/* Head-to-Head Data Display */}
        {selectedPlayer1 && selectedPlayer2 && selectedPlayer1 !== selectedPlayer2 && (
          <>
            {headToHeadLoading && (
              <div className="text-center py-12">
                <div className="athletic-spinner mx-auto mb-4"></div>
                <p className="athletic-body text-grey-300">Loading head-to-head data...</p>
              </div>
            )}

            {headToHeadError && (
              <div className="athletic-container text-center">
                <h2 className="athletic-heading text-2xl mb-4">No Games Played 😢</h2>
                <p className="athletic-body text-grey-300 mb-6">
                  Get ponging!
                </p>
                <button
                  onClick={() => window.location.reload()}
                  className="athletic-btn athletic-btn-primary"
                >
                  Retry
                </button>
              </div>
            )}

            {headToHeadData && <HeadToHeadDisplay data={headToHeadData} />}
          </>
        )}

        {/* No players selected state */}
        {(!selectedPlayer1 || !selectedPlayer2) && (
          <div className="text-center py-12">
            <div className="athletic-container max-w-md mx-auto">
              <h3 className="athletic-heading text-xl mb-4">Select Two Players</h3>
              <p className="athletic-body text-grey-300 mb-6">
                Choose two players from the dropdowns above to view their head-to-head statistics.
              </p>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default HeadToHeadPage;
