<script>
  import { onMount, onDestroy, tick } from 'svelte';
  import { lists } from './stores/lists.js';
  import ListCard from './components/ListCard.svelte';

  let newListName = '';
  let creatingList = false;
  let cursorVisible = true;
  let cursorInterval;
  let inputEl;

  onMount(() => { cursorInterval = setInterval(() => cursorVisible = !cursorVisible, 530); });
  onDestroy(() => clearInterval(cursorInterval));

  async function startCreating() {
    creatingList = true;
    await tick();
    inputEl?.focus();
  }

  function createList() {
    if (newListName.trim()) {
      lists.addList(newListName);
      newListName = '';
      creatingList = false;
    }
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') createList();
    if (e.key === 'Escape') { creatingList = false; newListName = ''; }
  }
</script>

<div class="app">
  <header>
    <div class="logo-line">
      <span class="bracket">[</span>
      <span class="logo-name">SISYPHUS</span>
      <span class="bracket">]</span>
      <span class="logo-ver"> v1.0.0</span>
    </div>
    <div class="sys-info">
      <span class="label">LISTS:</span>
      <span class="val">{$lists.length}</span>
    </div>
  </header>

  <div class="prompt-bar">
    <span class="ps1">root@sisyphus:~# </span>
    {#if creatingList}
<input
  class="cli-input"
  bind:this={inputEl}
  bind:value={newListName}
  placeholder="new_list_name"
  on:keydown={handleKeydown}
/>
    <span class="cursor" class:visible={cursorVisible}>█</span>
      <span class="hint"> [ENTER] confirm  [ESC] cancel</span>
    {:else}
<button class="inline-btn" on:click={startCreating}>+ NEW LIST</button>
      <span class="cursor" class:visible={cursorVisible}>█</span>
    {/if}
  </div>

  <main>
    {#if $lists.length === 0 && !creatingList}
      <div class="boot-screen">
        <pre class="ascii-boulder">{`
   ██████
  ████████
  ████████
   ██████`}</pre>
        <div class="boot-lines">
          <p class="line">SISYPHUS TASK MANAGER v1.0.0</p>
          <p class="line dim">Mounting storage ............. OK</p>
          <p class="line dim">Checking task lists .......... EMPTY</p>
          <p class="line warn">WARNING: No task lists found.</p>
          <p class="line">
            Press
<button class="key-btn" on:click={startCreating}>[F1] CREATE FIRST LIST</button>
          </p>
        </div>
      </div>
    {:else}
      <div class="lists-header">
        <span class="dir-marker">DIR</span>
        <span class="dir-path">/home/sisyphus/tasks/</span>
        <span class="dir-count">{$lists.length} entries</span>
      </div>
      <div class="grid">
        {#each $lists as list (list.id)}
          <ListCard {list} />
        {/each}
      </div>
    {/if}
  </main>

  <footer>
    <span>[SISYPHUS v1.0.0]</span>
    <span>[LISTS: {$lists.length}]</span>
    <span class="amber">[THE ROCK ROLLS ON]</span>
    <span class="right">[F1:NEW] [DBLCLK:EDIT]</span>
  </footer>
</div>

<style>
  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  /* HEADER */
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.5rem;
    border-bottom: 2px solid var(--border-bright);
    background: var(--bg);
    position: sticky;
    top: 0;
    z-index: 100;
  }
  .logo-line { font-size: 0.6rem; letter-spacing: 0.1em; }
  .bracket   { color: var(--muted); }
  .logo-name {
    color: var(--text-primary);
    text-shadow: 0 0 8px var(--accent-glow), 0 0 20px var(--accent-glow);
  }
  .logo-ver  { color: var(--muted); font-size: 0.45rem; }
  .sys-info  { font-size: 0.5rem; display: flex; gap: 0.4rem; }
  .label     { color: var(--muted); }
  .val       { color: var(--amber); text-shadow: 0 0 6px var(--amber-glow); }

  /* PROMPT BAR */
  .prompt-bar {
    display: flex;
    align-items: center;
    padding: 0.5rem 1.5rem;
    border-bottom: 1px solid var(--border);
    background: #010801;
    font-size: 0.5rem;
    gap: 0.2rem;
    flex-wrap: wrap;
  }
  .ps1 { color: var(--amber); text-shadow: 0 0 6px var(--amber-glow); white-space: nowrap; }
  .cli-input {
    background: transparent;
    border: none;
    outline: none;
    color: var(--text-primary);
    font-family: var(--font-pixel);
    font-size: 0.5rem;
    min-width: 160px;
    text-shadow: 0 0 6px var(--accent-glow);
  }
  .cursor { color: var(--text-primary); font-size: 0.55rem; visibility: hidden; }
  .cursor.visible { visibility: visible; }
  .hint { color: var(--muted); font-size: 0.42rem; }
  .inline-btn {
    background: var(--accent-subtle);
    border: 1px solid var(--border-bright);
    color: var(--accent);
    font-family: var(--font-pixel);
    font-size: 0.45rem;
    cursor: pointer;
    padding: 0.25rem 0.6rem;
    text-shadow: 0 0 6px var(--accent-glow);
    transition: all 0.1s;
  }
  .inline-btn:hover { background: var(--accent); color: #000; text-shadow: none; }

  /* MAIN */
  main {
    flex: 1;
    max-width: 960px;
    width: 100%;
    margin: 0 auto;
    padding: 1.5rem;
  }

  /* BOOT SCREEN */
  .boot-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    padding: 3rem 1rem;
  }
  .ascii-boulder {
    color: var(--muted);
    font-family: var(--font-pixel);
    font-size: 0.5rem;
    line-height: 1.4;
    text-align: center;
    text-shadow: 0 0 10px var(--accent-glow);
    animation: bob 2s ease-in-out infinite;
  }
  @keyframes bob {
    0%, 100% { transform: translateY(0); }
    50%       { transform: translateY(-6px); }
  }
  .boot-lines { display: flex; flex-direction: column; gap: 0.5rem; font-size: 0.5rem; }
  .line      { color: var(--text-secondary); }
  .line.dim  { color: var(--muted); }
  .line.warn { color: var(--amber); text-shadow: 0 0 8px var(--amber-glow); }
  .key-btn {
    background: var(--border-bright);
    color: #000;
    border: none;
    font-family: var(--font-pixel);
    font-size: 0.5rem;
    cursor: pointer;
    padding: 0.2rem 0.5rem;
  }
  .key-btn:hover { opacity: 0.8; }

  /* LISTS */
  .lists-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1rem;
    font-size: 0.45rem;
    border-left: 2px solid var(--border-bright);
    padding-left: 0.6rem;
  }
  .dir-marker { color: var(--amber); }
  .dir-path   { color: var(--text-secondary); }
  .dir-count  { color: var(--muted); margin-left: auto; }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
    gap: 1.25rem;
  }

  /* FOOTER */
  footer {
    display: flex;
    gap: 2px;
    padding: 0.4rem 0.6rem;
    border-top: 2px solid var(--border-bright);
    font-size: 0.42rem;
  }
  footer span {
    background: var(--muted);
    color: #000;
    padding: 0.2rem 0.5rem;
  }
  .amber { background: var(--amber) !important; }
  .right { margin-left: auto; background: var(--dim) !important; color: var(--muted) !important; }
</style>