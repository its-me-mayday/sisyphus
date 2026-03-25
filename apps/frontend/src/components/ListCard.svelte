<script>
  import { lists, listStats } from '../stores/lists.js';
  import TodoItem from './TodoItem.svelte';

  export let list;

  let editing = false;
  let editName = list.name;
  let expanded = true;
  let addingTodo = false;
  let newTodoText = '';

  $: stats = $listStats[list.id] ?? { total: 0, done: 0, pct: 0 };

  $: progressBar = (() => {
    const filled = Math.round(stats.pct / 10);
    return '█'.repeat(filled) + '░'.repeat(10 - filled);
  })();

  function saveListName() {
    if (editName.trim()) lists.updateList(list.id, editName);
    editing = false;
  }

  function handleListKeydown(e) {
    if (e.key === 'Enter') saveListName();
    if (e.key === 'Escape') { editing = false; editName = list.name; }
  }

  function addTodo() {
    if (newTodoText.trim()) {
      lists.addTodo(list.id, newTodoText);
      newTodoText = '';
      addingTodo = false;
    }
  }

  function handleTodoKeydown(e) {
    if (e.key === 'Enter') addTodo();
    if (e.key === 'Escape') { addingTodo = false; newTodoText = ''; }
  }

  function confirmDelete() {
    if (confirm(`DELETE: "${list.name}" ?`)) lists.deleteList(list.id);
  }
</script>

<div class="card">
  <!-- Title bar -->
  <div class="title-bar">
    <span class="tb-bracket">[</span>
    {#if editing}
      <input
        class="title-input"
        bind:value={editName}
        on:blur={saveListName}
        on:keydown={handleListKeydown}
      />
    {:else}
<span class="tb-name" on:dblclick={() => { editing = true; editName = list.name; }}>
  {list.name.toUpperCase()}
</span>
    {/if}
    <span class="tb-bracket">]</span>
    <span class="tb-spacer"></span>
    <div class="tb-btns">
      <button class="tb-btn" on:click={() => expanded = !expanded}>{expanded ? '▼' : '►'}</button>
      <button class="tb-btn" on:click={() => { editing = true; editName = list.name; }}>EDIT</button>
      <button class="tb-btn danger" on:click={confirmDelete}>DEL</button>
    </div>
  </div>

  <!-- Stats -->
  <div class="stats-row">
    <span class="stat-label">PROGRESS:</span>
    <span class="progress-bar">[<span class="bar-fill">{progressBar}</span>]</span>
    <span class="stat-pct" class:full={stats.pct === 100}>{stats.pct}%</span>
    <span class="stat-nums">{stats.done}/{stats.total}</span>
  </div>

  <!-- Body -->
  {#if expanded}
    <div class="body">
      {#each list.todos as todo (todo.id)}
        <TodoItem {todo} listId={list.id} />
      {/each}

      {#if list.todos.length === 0 && !addingTodo}
        <p class="empty">; no tasks — dir empty</p>
      {/if}

      {#if addingTodo}
        <div class="add-row">
          <span class="prompt">&gt; </span>
          <input
            class="add-input"
            bind:value={newTodoText}
            placeholder="task description..."
            on:keydown={handleTodoKeydown}
          />
          <button class="ok-btn" on:click={addTodo}>[OK]</button>
          <button class="cancel-btn" on:click={() => { addingTodo = false; newTodoText = ''; }}>[X]</button>
        </div>
      {:else}
        <button class="add-btn" on:click={() => addingTodo = true}>+ ADD TASK</button>
      {/if}
    </div>
  {/if}
</div>

<style>
  .card {
    background: var(--card-bg);
    border: 1px solid var(--border);
    border-top: 2px solid var(--border-bright);
    font-size: 0.5rem;
  }
  .card:hover { border-color: var(--border-bright); }

  /* Title bar */
  .title-bar {
    display: flex;
    align-items: center;
    padding: 0.5rem 0.6rem;
    background: var(--border-bright);
    color: #000;
    gap: 0.2rem;
  }
  .tb-bracket { color: rgba(0,0,0,0.4); font-size: 0.5rem; }
  .tb-name {
    font-size: 0.5rem;
    letter-spacing: 0.08em;
    color: #000;
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .title-input {
    background: transparent;
    border: none;
    border-bottom: 1px solid rgba(0,0,0,0.4);
    outline: none;
    font-family: var(--font-pixel);
    font-size: 0.5rem;
    color: #000;
    max-width: 160px;
    text-transform: uppercase;
  }
  .tb-spacer { flex: 1; }
  .tb-btns { display: flex; gap: 2px; }
  .tb-btn {
    background: rgba(0,0,0,0.2);
    border: 1px solid rgba(0,0,0,0.3);
    color: #000;
    font-family: var(--font-pixel);
    font-size: 0.38rem;
    cursor: pointer;
    padding: 0.2rem 0.4rem;
  }
  .tb-btn:hover { background: rgba(0,0,0,0.4); }
  .tb-btn.danger:hover { background: var(--danger); color: #fff; border-color: var(--danger); }

  /* Stats */
  .stats-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.45rem 0.7rem;
    border-bottom: 1px solid var(--border-faint);
    background: #010801;
  }
  .stat-label { color: var(--muted); font-size: 0.4rem; }
  .progress-bar { color: var(--muted); letter-spacing: -1px; }
  .bar-fill { color: var(--accent); text-shadow: 0 0 4px var(--accent-glow); }
  .stat-pct { color: var(--text-secondary); }
  .stat-pct.full { color: var(--accent); text-shadow: 0 0 8px var(--accent-glow); }
  .stat-nums { color: var(--muted); font-size: 0.4rem; margin-left: auto; }

  /* Body */
  .body { padding: 0.5rem 0.7rem 0.7rem; }

  .empty { color: var(--dim); font-size: 0.42rem; padding: 0.5rem 0; }

  .add-row {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    margin-top: 0.5rem;
    border: 1px solid var(--border-bright);
    padding: 0.3rem 0.4rem;
    background: var(--accent-subtle);
  }
  .prompt { color: var(--accent); }
  .add-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: var(--text-primary);
    font-family: var(--font-pixel);
    font-size: 0.45rem;
  }
  .ok-btn {
    background: var(--accent);
    border: none;
    color: #000;
    font-family: var(--font-pixel);
    font-size: 0.4rem;
    cursor: pointer;
    padding: 0.2rem 0.4rem;
  }
  .cancel-btn {
    background: none;
    border: 1px solid var(--danger);
    color: var(--danger);
    font-family: var(--font-pixel);
    font-size: 0.4rem;
    cursor: pointer;
    padding: 0.2rem 0.4rem;
  }

  .add-btn {
    margin-top: 0.6rem;
    background: none;
    border: 1px dashed var(--border);
    color: var(--muted);
    cursor: pointer;
    font-family: var(--font-pixel);
    font-size: 0.42rem;
    padding: 0.4rem 0.6rem;
    width: 100%;
    text-align: left;
    transition: all 0.1s;
  }
  .add-btn:hover {
    border-color: var(--border-bright);
    color: var(--accent);
    background: var(--accent-subtle);
  }
</style>