import { writable, derived } from 'svelte/store';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8090/api/v1';

function createListsStore() {
  const { subscribe, set, update } = writable([]);

  async function load() {
    try {
      const res = await fetch(`${BASE_URL}/lists`);
      const data = await res.json();
      set(data ?? []);
    } catch (e) {
      console.error('load lists:', e);
      set([]);
    }
  }

  return {
    subscribe,
    load,

    async addList(name) {
      const res = await fetch(`${BASE_URL}/lists`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      const list = await res.json();
      update(lists => [...lists, { ...list, todos: [] }]);
    },

    async updateList(id, name) {
      const res = await fetch(`${BASE_URL}/lists/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      const updated = await res.json();
      update(lists => lists.map(l => l.id === id ? { ...updated, todos: l.todos } : l));
    },

    async deleteList(id) {
      await fetch(`${BASE_URL}/lists/${id}`, { method: 'DELETE' });
      update(lists => lists.filter(l => l.id !== id));
    },

    async addTodo(listId, text, priority = 'medium', dueDate = '') {
      const res = await fetch(`${BASE_URL}/lists/${listId}/todos`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text, priority, due_date: dueDate }),
      });
      const todo = await res.json();
      update(lists => lists.map(l => l.id === listId
        ? { ...l, todos: [...l.todos, normalizeTodo(todo)] }
        : l
      ));
    },

    async toggleTodo(listId, todoId) {
      update(lists => {
        const list = lists.find(l => l.id === listId);
        const todo = list?.todos.find(t => t.id === todoId);
        if (!todo) return lists;

        fetch(`${BASE_URL}/todos/${todoId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            text:     todo.text,
            done:     !todo.done,
            priority: todo.priority ?? 'medium',
            due_date: todo.due_date ?? '',
          }),
        });

        return lists.map(l => l.id === listId ? {
          ...l,
          todos: l.todos.map(t => t.id === todoId ? { ...t, done: !t.done } : t)
        } : l);
      });
    },

    async updateTodo(listId, todoId, text) {
      update(lists => {
        const list = lists.find(l => l.id === listId);
        const todo = list?.todos.find(t => t.id === todoId);
        if (!todo) return lists;

        fetch(`${BASE_URL}/todos/${todoId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            text,
            done:     todo.done,
            priority: todo.priority ?? 'medium',
            due_date: todo.due_date ?? '',
          }),
        });

        return lists.map(l => l.id === listId ? {
          ...l,
          todos: l.todos.map(t => t.id === todoId ? { ...t, text } : t)
        } : l);
      });
    },

    async setTodoPriority(listId, todoId, priority) {
      update(lists => {
        const list = lists.find(l => l.id === listId);
        const todo = list?.todos.find(t => t.id === todoId);
        if (!todo) return lists;

        fetch(`${BASE_URL}/todos/${todoId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            text:     todo.text,
            done:     todo.done,
            priority,
            due_date: todo.due_date ?? '',
          }),
        });

        return lists.map(l => l.id === listId ? {
          ...l,
          todos: l.todos.map(t => t.id === todoId ? { ...t, priority } : t)
        } : l);
      });
    },

    async setTodoDueDate(listId, todoId, dueDate) {
      update(lists => {
        const list = lists.find(l => l.id === listId);
        const todo = list?.todos.find(t => t.id === todoId);
        if (!todo) return lists;

        fetch(`${BASE_URL}/todos/${todoId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            text:     todo.text,
            done:     todo.done,
            priority: todo.priority ?? 'medium',
            due_date: dueDate ?? '',
          }),
        });

        return lists.map(l => l.id === listId ? {
          ...l,
          todos: l.todos.map(t => t.id === todoId ? { ...t, due_date: dueDate } : t)
        } : l);
      });
    },

    async deleteTodo(listId, todoId) {
      await fetch(`${BASE_URL}/todos/${todoId}`, { method: 'DELETE' });
      update(lists => lists.map(l => l.id === listId ? {
        ...l,
        todos: l.todos.filter(t => t.id !== todoId)
      } : l));
    },

    reorderLists(newOrder) {
      update(() => newOrder);
    },
  };
}

// normalizza la risposta del backend — priority arriva come numero dal proto
function normalizeTodo(todo) {
  const PRIORITY_MAP = { 0: 'medium', 1: 'low', 2: 'medium', 3: 'high' };
  return {
    ...todo,
    priority: typeof todo.priority === 'number'
      ? PRIORITY_MAP[todo.priority] ?? 'medium'
      : todo.priority ?? 'medium',
  };
}

export const lists = createListsStore();

export const listStats = derived(lists, $lists =>
  Object.fromEntries($lists.map(l => {
    const total = l.todos?.length ?? 0;
    const done  = l.todos?.filter(t => t.done).length ?? 0;
    const pct   = total === 0 ? 0 : Math.round((done / total) * 100);
    return [l.id, { total, done, pending: total - done, pct }];
  }))
);