import React from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Trophy } from 'lucide-react';
import { api } from '../api';

const AchievementsPage: React.FC = () => {
  const { data: achievements, isLoading, error } = useQuery({
    queryKey: ['achievements'],
    queryFn: api.getAchievements,
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="text-center">
          <div className="athletic-spinner mx-auto mb-4"></div>
          <p className="athletic-body text-grey-300">Loading achievements...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="athletic-container max-w-md text-center">
          <h2 className="athletic-heading text-2xl mb-4">Error Loading Achievements</h2>
          <p className="athletic-body text-grey-300 mb-6">
            Unable to load achievements. Please try again.
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
      <header className="border-b border-grey-100 bg-white">
        <div className="max-w-6xl mx-auto px-6 py-6">
          <div className="flex items-center justify-between">
            <Link 
              to="/" 
              className="flex items-center gap-2 athletic-btn athletic-btn-secondary"
            >
              <ArrowLeft size={16} />
              Back to Leaderboard
            </Link>
            <div className="flex items-center gap-3">
              <Trophy className="text-lime-green" size={32} />
              <h1 className="athletic-display text-3xl">Achievements</h1>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="mb-8">
          <h2 className="athletic-display text-5xl mb-2">All Achievements</h2>
          <p className="athletic-label">UNLOCK MILESTONES BY PLAYING TABLE TENNIS</p>
        </div>

        {/* Achievements Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {achievements?.map((achievement) => (
            <div key={achievement.id} className="achievement-card">
              <div className="achievement-icon">
                <Trophy size={24} />
              </div>
              <div className="achievement-content">
                <h3 className="achievement-title">{achievement.title}</h3>
                <p className="achievement-description">{achievement.description}</p>
              </div>
            </div>
          ))}
        </div>

        {achievements?.length === 0 && (
          <div className="text-center py-12">
            <Trophy className="mx-auto mb-4 text-grey-300" size={48} />
            <h3 className="athletic-heading text-xl mb-2">No Achievements Available</h3>
            <p className="athletic-body text-grey-300">
              Achievements will appear here as they become available.
            </p>
          </div>
        )}
      </main>
    </div>
  );
};

export default AchievementsPage;