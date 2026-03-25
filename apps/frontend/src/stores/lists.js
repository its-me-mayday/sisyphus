import { writable, derived } from 'svelte/store';
import { v4 as uuidv4 } from 'uuid';

const STORAGE_KEY = 'sisyphus_lists';

function loadFromStorage() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function createListsStore() {
  const { subscribe, update } = writable(loadFromStorage());

  function persist(fn) {
    update(lists => {
      const next = fn(lists);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
      return next;
    });
  }

  return {
    subscribe,

    addList(name) {
      persist(lists => [...lists, {
        id: uuidv4(),
        name: name.trim(),
        createdAt: Date.now(),
        todos: [],
      }]);
    },

    updateList(id, name) {
      persist(lists => lists.map(l => l.id === id ? { ...l, name: name.trim() } : l));
    },

    deleteList(id) {
      persist(lists => lists.filter(l => l.id !== id));
    },

    reorderLists(newOrder) {
  persist(() => newOrder);
},

    addTodo(listId, text) {
      persist(lists => lists.map(l => l.id === listId ? {
        ...l,
        todos: [...l.todos, { id: uuidv4(), text: text.trim(), done: false, createdAt: Date.now() }]
      } : l));
    },

    toggleTodo(listId, todoId) {
      persist(lists => lists.map(l => l.id === listId ? {
        ...l,
        todos: l.todos.map(t => t.id === todoId ? { ...t, done: !t.done } : t)
      } : l));
    },

    updateTodo(listId, todoId, text) {
      persist(lists => lists.map(l => l.id === listId ? {
        ...l,
        todos: l.todos.map(t => t.id === todoId ? { ...t, text: text.trim() } : t)
      } : l));
    },

    deleteTodo(listId, todoId) {
      persist(lists => lists.map(l => l.id === listId ? {
        ...l,
        todos: l.todos.filter(t => t.id !== todoId)
      } : l));
    },
  };
}

export const lists = createListsStore();

export const listStats = derived(lists, $lists =>
  Object.fromEntries($lists.map(l => {
    const total = l.todos.length;
    const done  = l.todos.filter(t => t.done).length;
    const pct   = total === 0 ? 0 : Math.round((done / total) * 100);
    return [l.id, { total, done, pending: total - done, pct }];
  }))
);