<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useTitle } from '@vueuse/core';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { api } from '@/plugins/api';
import { setSpotifyProfile, spotifyProfileStore } from '@/stores/spotifyProfile';
import { selectedDashboardStore } from '@/stores/userStore';


const title = useTitle();
title.value = 'Tsuwari - Integrations';

const router = useRouter();

const selectedDashboard = useStore(selectedDashboardStore);

function spotifyRedirect() {
  window.location.replace(
    `https://accounts.spotify.com/authorize?` +
      new URLSearchParams({
        response_type: 'code',
        client_id: '47b25d8409384208873be26062de7243',
        scope: 'user-read-currently-playing',
        redirect_uri: `${window.location.origin}/dashboard/integrations`,
      }),
  );
}

const spotifyProfile = useStore(spotifyProfileStore);

async function fetchSpotifyProfile() {
  const { data } = await api(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/profile`);
  setSpotifyProfile(data);
}

onMounted(async () => {
  const params = new URLSearchParams(window.location.search);

  const code = params.get('code');
  if (code) {
    await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/token`, {
      code,
    });

    router.push('/dashboard/integrations');
  }

  if (!spotifyProfile.value) {
    fetchSpotifyProfile();
  }
});
</script>

<template>
  <div class="grid xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2 p-1">
    <div class="card rounded card-compact bg-base-200 drop-shadow-lg">
      <div class="card-body">
        <label class="label cursor-pointer mb-5">
          <h2 class="card-title outline-none">Donation Alerts</h2>
          <input
            type="checkbox"
            class="toggle"
            checked
          >
        </label>

        <div class="mb-5">
          <span class="label-text">Access token</span>
          <input
            type="text"
            placeholder="Donation Alerts Access token"
            class="rounded input input-sm input-bordered w-full max-w-xs mb-5"
          >

          <br>

          <span class="label-text">Refresh token</span>
          <input
            type="text"
            placeholder="Donation Alerts Refresh token"
            class="rounded input input-sm input-bordered w-full max-w-xs"
          >
        </div>

        <div class="card-actions">
          <button class="btn btn-outline btn-sm small-screen rounded">
            Generate Tokens
          </button>
          <button class="btn btn-secondary btn-sm small-screen rounded">
            Save
          </button>
        </div>
      </div>
    </div>

    <div class="card rounded card-compact bg-base-200 drop-shadow-lg">
      <div class="card-body">
        <label class="label cursor-pointer mb-5">
          <h2 class="card-title">Spotify</h2>
          <input
            type="checkbox"
            class="toggle"
            checked
          >
        </label>

        <div class="mb-5">
          <span v-if="spotifyProfile">Logged in as {{ spotifyProfile.display_name }}#{{ spotifyProfile.id }}</span>
          <span v-else>Not logged in</span>
        </div>

        <div class="card-actions">
          <button
            class="btn btn-outline btn-sm small-screen rounded"
            @click="spotifyRedirect"
          >
            Login
          </button>
        </div>
      </div>
    </div>

    <div class="card rounded card-compact bg-base-200 drop-shadow-lg">
      <div class="card-body">
        <label class="label cursor-pointer mb-5">
          <h2 class="card-title">QIWI</h2>
          <input
            type="checkbox"
            class="toggle"
            checked
          >
        </label>

        <div class="mb-5">
          <span class="label-text">Access token</span>
          <input
            type="text"
            placeholder="QIWI Access token"
            class="rounded input input-sm input-bordered w-full max-w-xs mb-5"
          >

          <br>

          <span class="label-text">Refresh token</span>
          <input
            type="text"
            placeholder="QIWI Refresh token"
            class="rounded input input-sm input-bordered w-full max-w-xs"
            disabled
          >
        </div>

        <div class="card-actions">
          <button class="btn btn-outline btn-sm small-screen rounded">
            Generate Tokens
          </button>
          <button class="btn btn-secondary btn-sm small-screen rounded">
            Save
          </button>
        </div>
      </div>
    </div>

    <div class="card rounded card-compact bg-base-200 drop-shadow-lg">
      <div class="card-body">
        <h2 class="card-title mb-5">
          Satont API
        </h2>

        <div class="mb-5">
          <label class="label cursor-pointer">
            <span class="label-text">Songs</span>
            <input
              type="checkbox"
              class="toggle"
              checked
            >
          </label>

          <label class="rounded input-group input-group-vertical mb-5">
            <span>VK ID</span>
            <input
              type="text"
              placeholder="VKontakte ID"
              class="rounded w-full input input-sm input-bordered"
            >
          </label>

          <label class="rounded input-group input-group-vertical mb-5">
            <span>Last FM</span>
            <input
              type="text"
              placeholder="Last FM ID"
              class="rounded w-full input input-sm input-bordered"
            >
          </label>

          <label class="rounded input-group input-group-vertical mb-5">
            <span>Twitch DJ</span>
            <input
              type="text"
              placeholder="Twitch DJ ID"
              class="rounded w-full input input-sm input-bordered"
            >
          </label>

          <label class="label cursor-pointer">
            <span class="label-text">FaceIT</span>
            <input
              type="checkbox"
              class="toggle"
              checked
            >
          </label>

          <label class="input-group input-group-vertical mb-5">
            <span>FaceIT</span>
            <input
              type="text"
              placeholder="Faceit nickname"
              class="rounded w-full input input-sm input-bordered"
            >
          </label>
        </div>

        <div class="card-actions justify-start">
          <button class="btn btn-secondary btn-sm w-full rounded">
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 426px) {
  .small-screen {
    width: 100%;
  }
}
</style>
