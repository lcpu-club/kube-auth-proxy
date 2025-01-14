<script setup lang="ts">
import { client } from '@/api/api';
import { ref } from 'vue';
import { RouterLink } from 'vue-router';

const tokens = ref((await (await client.get('/_/tokens')).json()).tokens as string[]);

const add_token = async () => {
  const token = await (await client.post('/_/tokens', {})).json();
  tokens.value.push(token.token);
};

const delete_token = async (token: string) => {
  await client.delete(`/_/tokens/${token}`);
  tokens.value = tokens.value.filter((t) => t !== token);
};

const kubeconfig_for = ref('');
const kubeconfig = ref('');

const show_kubeconfig = async (token: string) => {
  const kc = await (await client.get(`/_/tokens/${token}/kubeconfig`)).text();
  kubeconfig_for.value = token;
  kubeconfig.value = kc;
};
</script>

<template>
  <main>
    <h1>Tokens</h1>
    <p>
      <RouterLink :to="{ 'name': 'home' }">Back Home</RouterLink>
    </p>
    <p>Here are your tokens:</p>
    <ul>
      <li v-for="token in tokens" :key="token">
        <pre>{{ token }}</pre>
        <ul>
          <li><button @click="delete_token(token)">Delete</button></li>
          <li>
            <button @click="show_kubeconfig(token)">Show Kubeconfig</button>
          </li>
        </ul>
      </li>
    </ul>
    <p>
      <button @click="add_token()">Add Token</button>
    </p>
    <div v-if="kubeconfig">
      Kubeconfig for
      <pre style="display: inline;">{{ kubeconfig_for }}</pre> <br />
      <button @click="kubeconfig = '';">Hide</button> <br />
      <textarea rows="16" cols="80" disabled>{{ kubeconfig }}</textarea>
    </div>
  </main>
</template>
