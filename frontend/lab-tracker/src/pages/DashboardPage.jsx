import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { listLabWorks } from '../api/auth';
import { useAuth } from '../context/AuthContext';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const [labs, setLabs] = useState([]);
  const [meta, setMeta] = useState({ page: 1, totalPages: 0, total: 0, limit: 10 });
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadLabs = async (page = 1, value = search) => {
    setLoading(true);
    setError('');
    try {
      const { data } = await listLabWorks({ page, limit: 10, search: value });
      setLabs(data.items || []);
      setMeta({
        page: data.page || 1,
        totalPages: data.totalPages || 0,
        total: data.total || 0,
        limit: data.limit || 10,
      });
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось загрузить лабораторные работы');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLabs(1, '');
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const submitSearch = (e) => {
    e.preventDefault();
    loadLabs(1, search);
  };

  return (
    <div className="dash-wrap">
      <header className="dash-header">
        <Link className="dash-logo" to="/dashboard">
          <div className="auth-logo-icon">⬡</div>
          <div className="auth-logo-text">Хим<span>Лаб</span></div>
        </Link>
        <nav className="dash-nav">
          <Link to="/dashboard">Работы</Link>
          {user?.role === 'admin' && <Link to="/admin">Админка</Link>}
        </nav>
        <div className="dash-user">
          <span className="dash-username">{user?.username} · {user?.role || 'student'}</span>
          <button className="btn-logout" onClick={handleLogout}>Выйти</button>
        </div>
      </header>

      <main className="dash-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Лабораторные работы</h1>
          <p className="dash-subtitle">Каталог работ по химии, поиск и просмотр карточек</p>
        </div>

        <div className="dash-stats">
          <div className="stat-card">
            <div className="stat-value">{meta.total}</div>
            <div className="stat-label">Всего работ</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{meta.page}</div>
            <div className="stat-label">Текущая страница</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{labs.length}</div>
            <div className="stat-label">Показано на странице</div>
          </div>
        </div>

        <form className="toolbar" onSubmit={submitSearch}>
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Поиск по названию"
          />
          <button className="btn-secondary" type="submit">Найти</button>
        </form>

        <div className="dash-section-title">Список лабораторных работ</div>

        {error && <div className="auth-error">{error}</div>}
        {loading && <div className="empty-state">Загрузка...</div>}

        {!loading && labs.length === 0 && (
          <div className="empty-state">Записи не найдены</div>
        )}

        <div className="lab-list">
          {labs.map((lab, i) => (
            <Link className="lab-card" key={lab.id} to={`/labworks/${lab.id}`}>
              <div className="lab-index">Л/Р {String((meta.page - 1) * meta.limit + i + 1).padStart(2, '0')}</div>
              <div className="lab-info">
                <div className="lab-title">{lab.title}</div>
                {lab.goal && <div className="lab-status status-reviewed">{lab.goal}</div>}
              </div>
              <div className="lab-arrow">→</div>
            </Link>
          ))}
        </div>

        <div className="pagination">
          <button
            className="btn-secondary"
            disabled={meta.page <= 1}
            onClick={() => loadLabs(meta.page - 1)}
          >
            Назад
          </button>
          <span>{meta.page} / {meta.totalPages || 1}</span>
          <button
            className="btn-secondary"
            disabled={meta.totalPages === 0 || meta.page >= meta.totalPages}
            onClick={() => loadLabs(meta.page + 1)}
          >
            Вперёд
          </button>
        </div>
      </main>
    </div>
  );
}
