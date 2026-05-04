import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { getLabWork } from '../api/auth';
import { useAuth } from '../context/AuthContext';

export default function LabWorkPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  const [lab, setLab] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    getLabWork(id)
      .then(({ data }) => setLab(data))
      .catch((err) => setError(err.response?.data?.error?.message || 'Работа не найдена'));
  }, [id]);

  const handleLogout = () => {
    logout();
    navigate('/login');
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
          <span className="dash-username">{user?.username}</span>
          <button className="btn-logout" onClick={handleLogout}>Выйти</button>
        </div>
      </header>

      <main className="dash-main">
        {error && <div className="auth-error">{error}</div>}
        {!error && !lab && <div className="empty-state">Загрузка...</div>}

        {lab && (
          <>
            <div className="dash-greeting">
              <h1 className="dash-title">{lab.title}</h1>
              <p className="dash-subtitle">Карточка лабораторной работы</p>
            </div>

            <div className="details-grid">
              <Detail label="Описание" value={lab.description} />
              <Detail label="Цель" value={lab.goal} />
              <Detail label="Оборудование" value={lab.equipment} />
              <Detail label="Реактивы" value={lab.reagents} />
              <Detail label="Порядок выполнения" value={lab.procedure} />
              <Detail label="Файл" value={lab.filePath} />
            </div>
          </>
        )}
      </main>
    </div>
  );
}

function Detail({ label, value }) {
  if (!value) {
    return null;
  }

  return (
    <section className="detail-block">
      <div className="dash-section-title">{label}</div>
      <p>{value}</p>
    </section>
  );
}
