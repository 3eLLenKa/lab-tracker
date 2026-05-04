import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { login } from '../api/auth';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const navigate  = useNavigate();
  const location  = useLocation();
  const { saveUser } = useAuth();

  const [form,    setForm]    = useState({ username: '', password: '' });
  const [error,   setError]   = useState('');
  const [loading, setLoading] = useState(false);

  const registered = location.state?.registered;

  const handle = (e) =>
    setForm(f => ({ ...f, [e.target.name]: e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const { data } = await login(form);
      saveUser({ username: form.username }, data.token);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Неверный логин или пароль');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-wrap">
      <div className="auth-card">
        <div className="auth-logo">
          <div className="auth-logo-icon">⬡</div>
          <div className="auth-logo-text">Хим<span>Лаб</span></div>
        </div>

        <h1 className="auth-title">Вход в систему</h1>
        <p className="auth-subtitle">Трекер лабораторных работ</p>

        <form onSubmit={submit}>
          {registered && (
            <div className="auth-success">Аккаунт создан — войдите в систему</div>
          )}
          {error && <div className="auth-error">{error}</div>}

          <div className="field">
            <label>Имя пользователя</label>
            <input
              name="username"
              value={form.username}
              onChange={handle}
              autoComplete="username"
              required
            />
          </div>

          <div className="field">
            <label>Пароль</label>
            <input
              type="password"
              name="password"
              value={form.password}
              onChange={handle}
              autoComplete="current-password"
              required
            />
          </div>

          <button className="btn-primary" type="submit" disabled={loading}>
            <span>{loading ? 'Вход...' : 'Войти'}</span>
          </button>
        </form>

        <div className="auth-link">
          Нет аккаунта? <Link to="/register">Регистрация</Link>
        </div>

        <div className="chem-formula">H₂SO₄ · NaOH · KMnO₄</div>
      </div>
    </div>
  );
}