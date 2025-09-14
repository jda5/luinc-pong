import React, { useState, useEffect } from 'react';

interface AuthWrapperProps {
  children: React.ReactNode;
}

const AuthWrapper: React.FC<AuthWrapperProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(true);

  const correctPassword = 'game?';

  useEffect(() => {
    // Check if user is already authenticated
    const storedAuth = localStorage.getItem('table-tennis-auth');
    if (storedAuth === 'authenticated') {
      setIsAuthenticated(true);
    }
    setIsLoading(false);
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (password === correctPassword) {
      setIsAuthenticated(true);
      localStorage.setItem('table-tennis-auth', 'authenticated');
      setError('');
    } else {
      setError('Incorrect password. Try again!');
      setPassword('');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center px-4">
        <div className="athletic-container max-w-md w-full">
          <div className="text-center mb-8">
            <h1 className="athletic-display text-4xl mb-2">
              🏓 LUinc. Pong
            </h1>
            <p className="athletic-label">
              Enter Password
            </p>
          </div>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="password" className="athletic-label block mb-2">
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="athletic-input"
                placeholder="Enter password"
              />
            </div>
            
            {error && (
              <div className="text-red-600 text-sm text-center athletic-body">{error}</div>
            )}
            
            <div>
              <button
                type="submit"
                className="athletic-btn athletic-btn-primary w-full"
              >
                Enter
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      {children}
    </div>
  );
};

export default AuthWrapper;
