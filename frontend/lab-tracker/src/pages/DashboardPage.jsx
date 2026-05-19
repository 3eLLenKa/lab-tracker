import { useEffect, useMemo, useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import {
  getStudentProgress,
  listStudentAssignments,
  listTeacherSubmissions,
  setGrade,
  submitAssignment,
} from '../api/auth';
import { useAuth } from '../context/AuthContext';

const TEACHER_TABS = [
  { key: 'submitted', label: 'На проверке' },
  { key: 'revision', label: 'На доработке' },
  { key: 'reviewed', label: 'Проверенные' },
];

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
  const [progress, setProgress] = useState(null);
  const [activeAssignmentId, setActiveAssignmentId] = useState(null);
  const [reportForm, setReportForm] = useState({ textReport: '', filePath: '' });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');

  const loadAssignments = async () => {
    setLoading(true);
    try {
      const [{ data: assignmentsData }, { data: progressData }] = await Promise.all([
        listStudentAssignments(),
        getStudentProgress(),
      ]);
      setAssignments(assignmentsData.assignments || []);
      setProgress(progressData);
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

  const activeAssignment = assignments.find((item) => item.assignmentId === activeAssignmentId) || null;
  const activeState = activeAssignment?.submissionStatus || 'draft';
  const canEditActive = activeState === 'draft' || activeState === 'revision';
  const isReadOnlyActive = activeAssignment && !canEditActive;

  const submit = async (e) => {
    e.preventDefault();
    if (!activeAssignmentId || !canEditActive) {
      return;
    }

    setSaving(true);
    setError('');
    setNotice('');

    try {
      await submitAssignment({
        assignmentId: activeAssignmentId,
        textReport: reportForm.textReport,
        filePath: reportForm.filePath || null,
      });
      setNotice(activeState === 'revision' ? 'Исправленный отчёт отправлен повторно' : 'Отчёт отправлен на проверку');
      setActiveAssignmentId(null);
      setReportForm({ textReport: '', filePath: '' });
      await loadAssignments();
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось отправить отчёт');
    } finally {
      setSaving(false);
    }
  };

  const draftCount = assignments.filter((item) => !item.submissionStatus || item.submissionStatus === 'draft').length;
  const submittedCount = assignments.filter((item) => item.submissionStatus === 'submitted').length;
  const revisionCount = assignments.filter((item) => item.submissionStatus === 'revision').length;
  const reviewedCount = assignments.filter((item) => item.submissionStatus === 'reviewed').length;
  const completionRate = progress?.completionRate ?? 0;

  return (
    <div className="dash-wrap">
      <Header user={user} onLogout={handleLogout} />

      <main className="dash-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Мои лабораторные работы</h1>
          <p className="dash-subtitle">Список назначений, дедлайны, статусы и отправка отчётов</p>
        </div>

        <div className="dash-stats">
          <Stat value={assignments.length} label="Всего назначено" />
          <Stat value={draftCount} label="Готовятся" />
          <Stat value={submittedCount + revisionCount} label="В работе" />
          <Stat value={`${completionRate}%`} label={`Проверено: ${reviewedCount}`} />
        </div>

        {error && <div className="auth-error">{error}</div>}
        {notice && <div className="auth-success">{notice}</div>}
        {loading && <div className="empty-state">Загрузка...</div>}

        <div className="dash-section-title">Назначенные лабораторные</div>

        <div className="lab-list">
          {assignments.map((assignment) => {
            const status = assignment.submissionStatus || 'draft';
            const attemptText = assignment.attemptNumber ? `Попытка ${assignment.attemptNumber}` : 'Попытка 1';
            const editable = status === 'draft' || status === 'revision';

            return (
              <article className="student-card" key={assignment.assignmentId}>
                <div className="student-card-top">
                  <div>
                    <div className="lab-title">{assignment.title}</div>
                    <div className="student-meta">
                      Л/Р #{assignment.labWorkId}
                      {assignment.deadline && ` · дедлайн ${formatDate(assignment.deadline)}`}
                      {` · ${attemptText}`}
                    </div>
                  </div>
                  <span className={`lab-status ${statusClass(status)}`}>
                    {studentStatusLabel(status)}
                  </span>
                </div>

                <p className="student-description">{assignment.description}</p>

                {assignment.teacherComment && (
                  <div className={`teacher-note ${status === 'revision' ? 'teacher-note-revision' : ''}`}>
                    <strong>Комментарий преподавателя</strong>
                    <span>{assignment.teacherComment}</span>
                  </div>
                )}

                {(assignment.grade !== null && assignment.grade !== undefined) && (
                  <div className="grade-box">
                    <strong>Оценка: {assignment.grade}</strong>
                    {assignment.teacherComment && <span>{assignment.teacherComment}</span>}
                  </div>
                )}

                <div className="student-actions">
                  <Link className="btn-secondary" to={`/labworks/${assignment.labWorkId}`}>Карточка работы</Link>
                  <button className="btn-secondary" type="button" onClick={() => startSubmit(assignment)}>
                    {studentActionLabel(status)}
                  </button>
                </div>

                {activeAssignmentId === assignment.assignmentId && (
                  <form className="submission-form" onSubmit={submit}>
                    <div className="submission-state-line">
                      <span className={`lab-status ${statusClass(status)}`}>{studentStatusLabel(status)}</span>
                      {assignment.submittedAt && <span className="student-meta">Последняя отправка: {formatDate(assignment.submittedAt)}</span>}
                    </div>

                    {isReadOnlyActive && (
                      <div className="empty-state">
                        {status === 'reviewed'
                          ? 'Работа проверена и закрыта для повторной отправки.'
                          : 'Отчёт уже отправлен. Редактирование станет доступно только после возврата на доработку.'}
                      </div>
                    )}

                    <label className="admin-field">
                      <span>Текст отчёта</span>
                      <textarea
                        rows="6"
                        value={reportForm.textReport}
                        onChange={(e) => setReportForm((current) => ({ ...current, textReport: e.target.value }))}
                        disabled={!canEditActive || saving}
                        required
                      />
                    </label>

                    <label className="admin-field">
                      <span>Файл отчёта</span>
                      <input
                        type="file"
                        disabled={!canEditActive || saving}
                        onChange={(e) => {
                          const file = e.target.files?.[0];
                          setReportForm((current) => ({ ...current, filePath: file ? file.name : current.filePath }));
                        }}
                      />
                    </label>

                    {reportForm.filePath && <div className="file-chip">Файл: {reportForm.filePath}</div>}

                    <div className="form-actions">
                      {canEditActive && (
                        <button className="btn-primary" type="submit" disabled={saving}>
                          <span>{saving ? 'Отправка...' : status === 'revision' ? 'Отправить повторно' : 'Сдать на проверку'}</span>
                        </button>
                      )}
                      <button className="btn-secondary" type="button" onClick={() => setActiveAssignmentId(null)}>
                        Закрыть
                      </button>
                    </div>
                  </form>
                )}
              </article>
            );
          })}
        </div>
      </main>
    </div>
  );
}

function TeacherDashboard({ user, logout, navigate }) {
  const [submissions, setSubmissions] = useState([]);
  const [forms, setForms] = useState({});
  const [activeTab, setActiveTab] = useState('submitted');
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
      setForms((current) => ({
        ...Object.fromEntries(items.map((item) => [item.submissionId, {
          grade: item.grade ?? '',
          comment: item.teacherComment ?? '',
        }])),
        ...current,
      }));
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

  const counters = useMemo(() => ({
    submitted: submissions.filter((item) => item.status === 'submitted').length,
    revision: submissions.filter((item) => item.status === 'revision').length,
    reviewed: submissions.filter((item) => item.status === 'reviewed').length,
  }), [submissions]);

  const visibleSubmissions = useMemo(
    () => submissions.filter((item) => item.status === activeTab),
    [submissions, activeTab],
  );

  const updateForm = (submissionId, patch) => {
    setForms((current) => ({
      ...current,
      [submissionId]: { ...current[submissionId], ...patch },
    }));
  };

  const saveReview = async (submissionId, status) => {
    setSavingId(submissionId);
    setError('');
    setNotice('');

    try {
      const form = forms[submissionId] || { grade: '', comment: '' };
      await setGrade({
        submissionId,
        grade: Number(form.grade),
        comment: form.comment || null,
        status,
      });
      setNotice(status === 'revision' ? 'Работа отправлена на доработку' : 'Работа переведена в проверенные');
      await loadSubmissions();
      if (status === 'reviewed') {
        setActiveTab('submitted');
      }
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Не удалось сохранить результат проверки');
    } finally {
      setSavingId(null);
    }
  };

  return (
    <div className="dash-wrap">
      <Header user={user} onLogout={handleLogout} />

      <main className="dash-main teacher-main">
        <div className="dash-greeting">
          <h1 className="dash-title">Проверка отчётов</h1>
          <p className="dash-subtitle">Работы студентов своей группы с маршрутом проверки и доработки</p>
        </div>

        <div className="dash-stats">
          <Stat value={submissions.length} label="Всего отчётов" />
          <Stat value={counters.submitted} label="На проверке" />
          <Stat value={counters.revision} label="На доработке" />
          <Stat value={counters.reviewed} label="Проверено" />
        </div>

        {error && <div className="auth-error">{error}</div>}
        {notice && <div className="auth-success">{notice}</div>}
        {loading && <div className="empty-state">Загрузка...</div>}

        <div className="filter-tabs" role="tablist" aria-label="Фильтрация по статусу">
          {TEACHER_TABS.map((tab) => (
            <button
              key={tab.key}
              className={`filter-tab ${activeTab === tab.key ? 'filter-tab-active' : ''}`}
              onClick={() => setActiveTab(tab.key)}
              type="button"
            >
              <span>{tab.label}</span>
              <strong>{counters[tab.key]}</strong>
            </button>
          ))}
        </div>

        <div className="dash-section-title">Работы студентов</div>

        {visibleSubmissions.length === 0 && !loading && (
          <div className="empty-state">В этом разделе пока нет работ.</div>
        )}

        <div className="teacher-list">
          {visibleSubmissions.map((item) => {
            const isReviewed = item.status === 'reviewed';
            const form = forms[item.submissionId] || { grade: '', comment: '' };

            return (
              <article className="teacher-card" key={item.submissionId}>
                <div className="teacher-head">
                  <div>
                    <div className="lab-title">{item.labWorkTitle}</div>
                    <div className="teacher-meta">
                      {item.studentName} · {item.groupName}
                    </div>
                  </div>
                  <span className={`lab-status ${statusClass(item.status)}`}>{teacherStatusLabel(item.status)}</span>
                </div>

                <div className="teacher-details-grid">
                  <InfoLine label="Студент" value={item.studentName} />
                  <InfoLine label="Группа" value={item.groupName} />
                  <InfoLine label="Статус" value={teacherStatusLabel(item.status)} />
                  <InfoLine label="Попытка" value={`#${item.attemptNumber || 1}`} />
                  <InfoLine label="Дата сдачи" value={item.submittedAt ? formatDate(item.submittedAt) : 'Не отправлено'} />
                  <InfoLine label="Дедлайн" value={item.deadline ? formatDate(item.deadline) : 'Не указан'} />
                </div>

                <div className="teacher-report">{item.textReport}</div>
                {item.filePath && <div className="file-chip">Файл: {item.filePath}</div>}

                {item.teacherComment && (
                  <div className="teacher-note">
                    <strong>Последний комментарий</strong>
                    <span>{item.teacherComment}</span>
                  </div>
                )}

                <div className="grade-form-grid">
                  <label className="admin-field">
                    <span>Оценка (0-100)</span>
                    <input
                      type="number"
                      min="0"
                      max="100"
                      value={form.grade}
                      disabled={isReviewed || savingId === item.submissionId}
                      onChange={(e) => updateForm(item.submissionId, { grade: e.target.value })}
                    />
                  </label>
                  <label className="admin-field">
                    <span>Комментарий преподавателя</span>
                    <textarea
                      rows="4"
                      value={form.comment}
                      disabled={isReviewed || savingId === item.submissionId}
                      onChange={(e) => updateForm(item.submissionId, { comment: e.target.value })}
                    />
                  </label>
                </div>

                <div className="form-actions">
                  {!isReviewed && (
                    <>
                      <button
                        className="btn-secondary"
                        type="button"
                        disabled={savingId === item.submissionId}
                        onClick={() => saveReview(item.submissionId, 'revision')}
                      >
                        Отправить на доработку
                      </button>
                      <button
                        className="btn-primary"
                        type="button"
                        disabled={savingId === item.submissionId}
                        onClick={() => saveReview(item.submissionId, 'reviewed')}
                      >
                        <span>{savingId === item.submissionId ? 'Сохранение...' : 'Завершить проверку'}</span>
                      </button>
                    </>
                  )}
                </div>
              </article>
            );
          })}
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

function InfoLine({ label, value }) {
  return (
    <div className="info-line">
      <span>{label}</span>
      <strong>{value}</strong>
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
  if (status === 'revision') return 'status-revision';
  if (status === 'reviewed') return 'status-reviewed';
  return 'status-draft';
}

function studentStatusLabel(status) {
  if (status === 'submitted') return 'На проверке';
  if (status === 'revision') return 'На доработке';
  if (status === 'reviewed') return 'Проверено';
  return 'Черновик';
}

function teacherStatusLabel(status) {
  if (status === 'submitted') return 'На проверке';
  if (status === 'revision') return 'На доработке';
  if (status === 'reviewed') return 'Проверено';
  return 'Черновик';
}

function studentActionLabel(status) {
  if (status === 'submitted') return 'Посмотреть отчёт';
  if (status === 'reviewed') return 'Посмотреть результат';
  if (status === 'revision') return 'Доработать отчёт';
  return 'Сдать работу';
}
