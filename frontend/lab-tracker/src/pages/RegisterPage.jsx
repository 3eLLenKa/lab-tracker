import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { listGroups, register } from '../api/auth';

export default function RegisterPage() {
  const navigate = useNavigate();

  const [groups, setGroups] = useState([]);
  const [form, setForm] = useState({
    username: '',
    full_name: '',
    password: '',
    password2: '',
    group_id: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    listGroups()
      .then(r => setGroups(r.data.groups || []))
      .catch(() => {});
  }, []);

  const handle = (e) =>
    setForm(f => ({ ...f, [e.target.name]: e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    setError('');

    if (form.password !== form.password2) {
      return setError('Пароли не совпадают');
    }

    setLoading(true);
    try {
      await register({
        username: form.username,
        password: form.password,
        fullName: form.full_name,
        groupId: form.group_id ? parseInt(form.group_id) : null,
      });
      navigate('/login', { state: { registered: true } });
    } catch (err) {
      setError(err.response?.data?.error?.message || 'Ошибка регистрации.');
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

        <h1 className="auth-title">Регистрация</h1>
        <p className="auth-subtitle">Создать аккаунт студента</p>

        <form onSubmit={submit}>
          {error && <div className="auth-error">{error}</div>}

          <div className="field">
            <label>Полное имя</label>
            <input
              name="full_name"
              value={form.full_name}
              onChange={handle}
              placeholder="Иванов Иван Иванович"
              required
            />
          </div>

          <div className="field">
            <label>Имя пользователя</label>
            <input
              name="username"
              value={form.username}
              onChange={handle}
              placeholder="ivanov_ii"
              autoComplete="username"
              required
            />
          </div>

          <div className="field">
            <label>Учебная группа</label>
            <select name="group_id" value={form.group_id} onChange={handle}>
              <option value="">— Выбрать группу —</option>
              {groups.map(g => (
                <option key={g.id} value={g.id}>{g.name}</option>
              ))}
            </select>
          </div>

          <div className="field-row">
            <div className="field">
              <label>Пароль</label>
              <input
                type="password"
                name="password"
                value={form.password}
                onChange={handle}
                autoComplete="new-password"
                required
              />
            </div>
            <div className="field">
              <label>Повторить</label>
              <input
                type="password"
                name="password2"
                value={form.password2}
                onChange={handle}
                autoComplete="new-password"
                required
              />
            </div>
          </div>

          <button className="btn-primary" type="submit" disabled={loading}>
            <span>{loading ? 'Регистрация...' : 'Создать аккаунт'}</span>
          </button>
        </form>

        <div className="auth-link">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </div>

        <div className="chem-formula">NaCl · CuSO₄ · Ca(OH)₂</div>
      </div>
    </div>
  );
}
