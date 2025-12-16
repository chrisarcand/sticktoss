<script>
  import { Router, Route } from 'svelte-routing';
  import Login from './routes/Login.svelte';
  import Signup from './routes/Signup.svelte';
  import Dashboard from './routes/Dashboard.svelte';
  import Group from './routes/Group.svelte';
</script>

<Router>
  <Route path="/login" component={Login} />
  <Route path="/signup" component={Signup} />
  <Route path="/group/:id" component={Group} />
  <Route path="/" component={Dashboard} />
</Router>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background-color: #f5f5f5;
  }

  :global(*) {
    box-sizing: border-box;
  }
</style>
