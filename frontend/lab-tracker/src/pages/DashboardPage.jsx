import { useEffect, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import {
  listStudentAssignments,
  listTeacherSubmissions,
  setGrade,
  submitAssignment,
} from '../api/auth';
import { useAuth } from '../context/AuthContext';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  if (user?.role === 'admin') {
    return <Navigate to="/admin" replace />;
  }

  return user?.role === 'teacher'
    ? <TeacherDashboard user={user} logout={logout} navigate={navigate} />
    : <StudentDashboard user={user} logout={logout} navigate={navigate} />;
}

function StudentDashboard({ user, logout, navigate }) {
  const [assignments, setAssignments] = useState([]);
  const [activeAssignmentId, setActiveAssignmentId] = useState(null);
  const [reportForm, setReportForm] = useState({ textReport: '', filePath: '' });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');

  const loadAssignments = async () => {
    setLoading(true);
    try {
      const { data } = await listStudentAssignments();
      setAssignments(data.assignments || []);
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось загрузить назначения');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAssignments();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const startSubmit = (assignment) => {
    setActiveAssignmentId(assignment.assignmentId);
    setNotice('');
    setError('');
    setReportForm({
      textReport: assignment.textReport || '',
      filePath: assignment.filePath || '',
    });
  };

  const submit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setNotice('');

    try {
      await submitAssignment({
        assignmentId: activeAssignmentId,
        textReport: reportForm.textReport,
        filePath: reportForm.filePath || null,
      });
      setNotice('Отчёт отправлен на проверку');
      setActiveAssignmentId(null);
      setReportForm({ textReport: '', filePath: '' });
      await loadAssignments();
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось отправить отчёт');
    } finally {
      setSaving(false);
    }
  };

  const pendingCount = assignments.filter((item) => !item.submissionStatus || item.submissionStatus === 'pending').length;
  const checkedCount = assignments.filter((item) => item.submissionStatus === 'checked').length;

  return (
    <div className="dash-wrap">
      <Header user={user} onLogout={handleLogout} />

      <main className="dash-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Мои лабораторные работы</h1>
          <p className="dash-subtitle">Список назначений, дедлайны и сдача отчётов по химии</p>
        </div>

        <div className="dash-stats">
          <Stat value={assignments.length} label="Всего назначено" />
          <Stat value={pendingCount} label="Ожидают сдачи" />
          <Stat value={checkedCount} label="Проверено" />
        </div>

        {error && <div className="auth-error">{error}</div>}
        {notice && <div className="auth-success">{notice}</div>}
        {loading && <div className="empty-state">Загрузка...</div>}

        <div className="dash-section-title">Назначенные лабораторные</div>

        <div className="lab-list">
          {assignments.map((assignment) => (
            <article className="student-card" key={assignment.assignmentId}>
              <div className="student-card-top">
                <div>
                  <div className="lab-title">{assignment.title}</div>
                  <div className="student-meta">
                    Л/Р #{assignment.labWorkId} {assignment.deadline && `· дедлайн ${formatDate(assignment.deadline)}`}
                  </div>
                </div>
                <span className={`lab-status ${statusClass(assignment.submissionStatus)}`}>
                  {studentStatusLabel(assignment.submissionStatus)}
                </span>
              </div>

              <p className="student-description">{assignment.description}</p>

              {(assignment.grade !== null && assignment.grade !== undefined) && (
                <div className="grade-box">
                  <strong>Оценка: {assignment.grade}</strong>
                  {assignment.teacherComment && <span>{assignment.teacherComment}</span>}
                </div>
              )}

              <div className="student-actions">
                <Link className="btn-secondary" to={`/labworks/${assignment.labWorkId}`}>Карточка работы</Link>
                <button className="btn-secondary" type="button" onClick={() => startSubmit(assignment)}>
                  {assignment.submissionStatus === 'submitted' || assignment.submissionStatus === 'checked' ? 'Посмотреть отчёт' : 'Сдать работу'}
                </button>
              </div>

              {activeAssignmentId === assignment.assignmentId && (
                <form className="submission-form" onSubmit={submit}>
                  <label className="admin-field">
                    <span>Текст отчёта</span>
                    <textarea
                      rows="6"
                      value={reportForm.textReport}
                      onChange={(e) => setReportForm((current) => ({ ...current, textReport: e.target.value }))}
                      required
                    />
                  </label>

                  <label className="admin-field">
                    <span>Файл отчёта</span>
                    <input
                      type="file"
                      onChange={(e) => {
                        const file = e.target.files?.[0];
                        setReportForm((current) => ({ ...current, filePath: file ? file.name : '' }));
                      }}
                    />
                  </label>

                  {reportForm.filePath && <div className="file-chip">Выбран файл: {reportForm.filePath}</div>}

                  <div className="form-actions">
                    <button className="btn-primary" type="submit" disabled={saving}>
                      <span>{saving ? 'Отправка...' : 'Сдать на проверку'}</span>
                    </button>
                    <button className="btn-secondary" type="button" onClick={() => setActiveAssignmentId(null)}>
                      Закрыть
                    </button>
                  </div>
                </form>
              )}
            </article>
          ))}
        </div>
      </main>
    </div>
  );
}

function TeacherDashboard({ user, logout, navigate }) {
  const [submissions, setSubmissions] = useState([]);
  const [forms, setForms] = useState({});
  const [loading, setLoading] = useState(true);
  const [savingId, setSavingId] = useState(null);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');

  const loadSubmissions = async () => {
    setLoading(true);
    try {
      const { data } = await listTeacherSubmissions();
      const items = data.submissions || [];
      setSubmissions(items);
      setForms(Object.fromEntries(items.map((item) => [item.submissionId, {
        grade: item.grade ?? '',
        comment: item.teacherComment ?? '',
      }])));
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось загрузить отчёты группы');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSubmissions();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const saveGrade = async (submissionId) => {
    setSavingId(submissionId);
    setError('');
    setNotice('');
    try {
      const form = forms[submissionId];
      await setGrade({
        submissionId,
        grade: Number(form.grade),
        comment: form.comment,
      });
      setNotice('Оценка сохранена');
      await loadSubmissions();
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось сохранить оценку');
    } finally {
      setSavingId(null);
    }
  };

  const submittedCount = submissions.filter((item) => item.status === 'submitted').length;
  const checkedCount = submissions.filter((item) => item.status === 'checked').length;

  return (
    <div className="dash-wrap">
      <Header user={user} onLogout={handleLogout} />

      <main className="dash-main teacher-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Отчёты студентов</h1>
          <p className="dash-subtitle">Проверка лабораторных работ своей группы</p>
        </div>

        <div className="dash-stats">
          <Stat value={submissions.length} label="Всего отчётов" />
          <Stat value={submittedCount} label="Ждут проверки" />
          <Stat value={checkedCount} label="Уже оценены" />
        </div>

        {error && <div className="auth-error">{error}</div>}
        {notice && <div className="auth-success">{notice}</div>}
        {loading && <div className="empty-state">Загрузка...</div>}

        <div className="dash-section-title">Работы студентов</div>

        <div className="teacher-list">
          {submissions.map((item) => (
            <article className="teacher-card" key={item.submissionId}>
              <div className="teacher-head">
                <div>
                  <div className="lab-title">{item.labWorkTitle}</div>
                  <div className="teacher-meta">
                    {item.studentName} · {item.groupName} {item.submittedAt && `· ${formatDate(item.submittedAt)}`}
                  </div>
                </div>
                <span className={`lab-status ${statusClass(item.status)}`}>{teacherStatusLabel(item.status)}</span>
              </div>

              <div className="teacher-report">{item.textReport}</div>
              {item.filePath && <div className="file-chip">Файл: {item.filePath}</div>}

              <div className="grade-form-grid">
                <label className="admin-field">
                  <span>Оценка (0-100)</span>
                  <input
                    type="number"
                    min="0"
                    max="100"
                    value={forms[item.submissionId]?.grade ?? ''}
                    onChange={(e) => setForms((current) => ({
                      ...current,
                      [item.submissionId]: { ...current[item.submissionId], grade: e.target.value },
                    }))}
                  />
                </label>
                <label className="admin-field">
                  <span>Комментарий</span>
                  <textarea
                    rows="4"
                    value={forms[item.submissionId]?.comment ?? ''}
                    onChange={(e) => setForms((current) => ({
                      ...current,
                      [item.submissionId]: { ...current[item.submissionId], comment: e.target.value },
                    }))}
                  />
                </label>
              </div>

              <div className="form-actions">
                <button className="btn-primary" type="button" disabled={savingId === item.submissionId} onClick={() => saveGrade(item.submissionId)}>
                  <span>{savingId === item.submissionId ? 'Сохранение...' : 'Выставить оценку'}</span>
                </button>
              </div>
            </article>
          ))}
        </div>
      </main>
    </div>
  );
}

function Header({ user, onLogout }) {
  return (
    <header className="dash-header">
      <Link className="dash-logo" to="/dashboard">
        <div className="auth-logo-icon">⬡</div>
        <div className="auth-logo-text">Хим<span>Лаб</span></div>
      </Link>
      <nav className="dash-nav">
        <Link to="/dashboard">Кабинет</Link>
        {user?.role === 'admin' && <Link to="/admin">Админка</Link>}
      </nav>
      <div className="dash-user">
        <span className="dash-username">{user?.username} · {user?.role}</span>
        <button className="btn-logout" onClick={onLogout}>Выйти</button>
      </div>
    </header>
  );
}

function Stat({ value, label }) {
  return (
    <div className="stat-card">
      <div className="stat-value">{value}</div>
      <div className="stat-label">{label}</div>
    </div>
  );
}

function formatDate(value) {
  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value));
}

function statusClass(status) {
  if (status === 'submitted') return 'status-submitted';
  if (status === 'checked') return 'status-reviewed';
  return 'status-pending';
}

function studentStatusLabel(status) {
  if (status === 'submitted') return 'На проверке';
  if (status === 'checked') return 'Проверено';
  return 'Не сдано';
}

function teacherStatusLabel(status) {
  if (status === 'checked') return 'Оценено';
  if (status === 'submitted') return 'Ждёт проверки';
  return 'Черновик';
}
