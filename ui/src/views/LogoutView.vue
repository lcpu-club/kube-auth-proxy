<template>
    <div class="logout">
        <h1>Logging Out...</h1>
    </div>
</template>

<script setup lang="ts">
import { clearToken } from '@/api/token';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

onMounted(() => {
    clearToken();
    router.push({ 'name': 'home' });
});
</script>
