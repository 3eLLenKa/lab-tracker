import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { useAuth } from './context/AuthContext';
import LoginPage     from './pages/LoginPage';
import RegisterPage  from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import LabWorkPage   from './pages/LabWorkPage';
import AdminPage     from './pages/AdminPage';

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  );
}

function AppRoutes() {
  const { user } = useAuth();

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={user ? <Navigate to="/dashboard" replace /> : <LoginPage />} />
        <Route path="/register" element={user ? <Navigate to="/dashboard" replace /> : <RegisterPage />} />
        <Route path="/dashboard" element={<RequireAuth user={user}><DashboardPage /></RequireAuth>} />
        <Route path="/labworks/:id" element={<RequireAuth user={user}><LabWorkPage /></RequireAuth>} />
        <Route path="/admin" element={<RequireAdmin user={user}><AdminPage /></RequireAdmin>} />
        <Route path="*" element={<Navigate to={user ? '/dashboard' : '/login'} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

function RequireAuth({ user, children }) {
  return user ? children : <Navigate to="/login" replace />;
}

function RequireAdmin({ user, children }) {
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return user.role === 'admin' ? children : <Navigate to="/dashboard" replace />;
}

export default App;
