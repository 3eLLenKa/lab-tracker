import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  createLabWork,
  deleteLabWork,
  exportAdminCsv,
  getAdminStats,
  listLabWorks,
  updateLabWork,
} from '../api/auth';
import { useAuth } from '../context/AuthContext';

const EMPTY_FORM = {
  title: '',
  description: '',
  goal: '',
  equipment: '',
  reagents: '',
  procedure: '',
  filePath: '',
};

export default function AdminPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const [labs, setLabs] = useState([]);
  const [stats, setStats] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');

  async function loadLabs() {
    setLoading(true);
    try {
      const { data } = await listLabWorks({ page: 1, limit: 20 });
      setLabs(data.items || []);
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось загрузить записи');
    } finally {
      setLoading(false);
    }
  }

  async function loadStats() {
    try {
      const { data } = await getAdminStats();
      setStats(data);
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось загрузить статистику');
    }
  }

  useEffect(() => {
    loadLabs();
    loadStats();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((current) => ({ ...current, [name]: value }));
  };

  const resetForm = () => {
    setForm(EMPTY_FORM);
    setEditingId(null);
    setError('');
  };

  const startEdit = (lab) => {
    setEditingId(lab.id);
    setNotice('');
    setError('');
    setForm({
      title: lab.title || '',
      description: lab.description || '',
      goal: lab.goal || '',
      equipment: lab.equipment || '',
      reagents: lab.reagents || '',
      procedure: lab.procedure || '',
      filePath: lab.filePath || '',
    });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const submit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setNotice('');

    const payload = {
      title: form.title.trim(),
      description: form.description.trim(),
      goal: form.goal.trim(),
      equipment: form.equipment.trim(),
      reagents: form.reagents.trim(),
      procedure: form.procedure.trim(),
      filePath: form.filePath.trim() ? form.filePath.trim() : null,
    };

    try {
      if (editingId) {
        await updateLabWork(editingId, payload);
        setNotice('Запись обновлена');
      } else {
        await createLabWork(payload);
        setNotice('Запись добавлена');
      }
      resetForm();
      await loadLabs();
      await loadStats();
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось сохранить запись');
    } finally {
      setSaving(false);
    }
  };

  const remove = async (id) => {
    const confirmed = window.confirm('Удалить лабораторную работу?');
    if (!confirmed) {
      return;
    }

    setError('');
    setNotice('');
    try {
      await deleteLabWork(id);
      if (editingId === id) {
        resetForm();
      }
      setNotice('Запись удалена');
      await loadLabs();
      await loadStats();
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось удалить запись');
    }
  };

  const exportCsv = async () => {
    setError('');
    setNotice('');

    try {
      const { data } = await exportAdminCsv();
      const url = window.URL.createObjectURL(data);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'lab-tracker-report.csv';
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      setNotice('CSV-отчёт сформирован');
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось сформировать CSV');
    }
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
          <Link to="/admin">Админка</Link>
        </nav>
        <div className="dash-user">
          <span className="dash-username">{user?.username} · admin</span>
          <button className="btn-logout" onClick={handleLogout}>Выйти</button>
        </div>
      </header>

      <main className="dash-main admin-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Управление лабораторными</h1>
          <p className="dash-subtitle">CRUD, статистика и экспорт журнала сдач</p>
        </div>

        <section className="admin-panel admin-stats-panel">
          <div className="admin-panel-head">
            <div className="dash-section-title">Статистика проекта</div>
            <button className="btn-secondary" type="button" onClick={exportCsv}>
              Экспорт CSV
            </button>
          </div>

          <div className="admin-stats-grid">
            <AdminStat value={stats?.labWorksTotal ?? 0} label="Лабораторных" />
            <AdminStat value={stats?.studentsTotal ?? 0} label="Студентов" />
            <AdminStat value={stats?.teachersTotal ?? 0} label="Преподавателей" />
            <AdminStat value={stats?.submissionsTotal ?? 0} label="Сдач" />
            <AdminStat value={stats?.submittedCount ?? 0} label="На проверке" />
            <AdminStat value={stats?.revisionCount ?? 0} label="На доработке" />
            <AdminStat value={stats?.reviewedCount ?? 0} label="Проверено" />
            <AdminStat value={formatAverage(stats?.averageGrade)} label="Средний балл" />
          </div>
        </section>

        <div className="admin-grid">
          <section className="admin-panel">
            <div className="dash-section-title">
              {editingId ? `Редактирование #${editingId}` : 'Новая запись'}
            </div>

            {editingId && (
              <div className="edit-banner">
                Сейчас редактируется: <strong>{form.title || 'без названия'}</strong>
              </div>
            )}

            {error && <div className="auth-error">{error}</div>}
            {notice && <div className="auth-success">{notice}</div>}

            <form className="admin-form" onSubmit={submit}>
              <Field label="Название">
                <input name="title" value={form.title} onChange={handleChange} required />
              </Field>
              <Field label="Описание">
                <textarea name="description" value={form.description} onChange={handleChange} rows="4" />
              </Field>
              <Field label="Цель">
                <textarea name="goal" value={form.goal} onChange={handleChange} rows="3" />
              </Field>
              <Field label="Оборудование">
                <textarea name="equipment" value={form.equipment} onChange={handleChange} rows="3" />
              </Field>
              <Field label="Реактивы">
                <textarea name="reagents" value={form.reagents} onChange={handleChange} rows="3" />
              </Field>
              <Field label="Порядок выполнения">
                <textarea name="procedure" value={form.procedure} onChange={handleChange} rows="5" />
              </Field>
              <Field label="Путь к файлу">
                <input name="filePath" value={form.filePath} onChange={handleChange} />
              </Field>

              <div className="form-actions">
                <button className="btn-primary" type="submit" disabled={saving}>
                  <span>{saving ? 'Сохранение...' : editingId ? 'Сохранить изменения' : 'Добавить запись'}</span>
                </button>
                {editingId && (
                  <button className="btn-secondary" type="button" onClick={resetForm}>
                    Отменить
                  </button>
                )}
              </div>
            </form>
          </section>

          <section className="admin-panel">
            <div className="dash-section-title">Список записей</div>
            {loading && <div className="empty-state">Загрузка...</div>}
            {!loading && labs.length === 0 && <div className="empty-state">Записей пока нет</div>}

            <div className="admin-list">
              {labs.map((lab) => (
                <article className="admin-item" key={lab.id}>
                  <div className="admin-item-main">
                  <div className="admin-item-title">{lab.title}</div>
                    <div className="admin-item-meta">Лабораторная работа #{lab.id}</div>
                    {lab.description && <p className="admin-item-text">{lab.description}</p>}
                  </div>
                  <div className="admin-item-actions">
                    <button className="btn-secondary" type="button" onClick={() => startEdit(lab)}>
                      Изменить
                    </button>
                    <button className="btn-danger" type="button" onClick={() => remove(lab.id)}>
                      Удалить
                    </button>
                  </div>
                </article>
              ))}
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}

function Field({ label, children }) {
  return (
    <label className="admin-field">
      <span>{label}</span>
      {children}
    </label>
  );
}

function AdminStat({ value, label }) {
  return (
    <div className="admin-stat">
      <strong>{value}</strong>
      <span>{label}</span>
    </div>
  );
}

function formatAverage(value) {
  if (value === null || value === undefined) {
    return '—';
  }
  return Number(value).toFixed(1);
}
