<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { api } from '@/plugins/api';
import { setSpotifyProfile, spotifyProfileStore } from '@/stores/spotifyProfile';
import { selectedDashboardStore } from '@/stores/userStore';

const router = useRouter();
const spotifyIntegration = ref<Partial<ChannelIntegration>>({
  enabled: true,
});
const selectedDashboard = useStore(selectedDashboardStore);

selectedDashboardStore.subscribe(d => {
  api(`/v1/channels/${d.channelId}/integrations/spotify`).then(async (r) => {
    console.log('spotify', r.data);
    spotifyIntegration.value = r.data;
  });
});

function spotifyRedirect() {
  window.location.replace(
    `https://accounts.spotify.com/authorize?` +
      new URLSearchParams({
        response_type: 'code',
        client_id: '47b25d8409384208873be26062de7243',
        scope: 'user-read-currently-playing',
        redirect_uri: `${window.location.origin}/dashboard/integrations/spotify`,
      }),
  );
}

const spotifyProfile = useStore(spotifyProfileStore);

async function fetchSpotifyProfile() {
  const { data } = await api(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/profile`);
  setSpotifyProfile(data);
}

async function patch() {
  const { data } = await api.patch(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify`, {
    enabled: spotifyIntegration.value.enabled,
  });

  spotifyIntegration.value = data;
}

onMounted(async () => {
  const route = router.currentRoute.value;
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');

  if (route.params.integration === 'spotify' && code) {
    if (code) {
      await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/spotify/token`, {
        code,
      });

      router.push('/dashboard/integrations');
    }
  }

  if (!spotifyProfile.value) {
    fetchSpotifyProfile();
  }
});
</script>

<template>
  <div class="card rounded card-compact bg-base-200 drop-shadow-lg p-4">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title">
          Spotify
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="spotifyIntegration.enabled"
          class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow-sm"
          type="checkbox"
          role="switch"
          @change="patch"
        >
      </div>
    </div>

    <div class="mb-5">
      <div
        v-if="spotifyProfile"
      >
        <div class="flex justify-center mb-3">
          <img
            v-if="spotifyProfile.images"
            :src="spotifyProfile.images[0].url"
            class="rounded-full w-32 ring-2 ring-white"
            alt="Avatar"
          >
        </div>
        {{ spotifyProfile.display_name }}#{{ spotifyProfile.id }}
      </div>
      <div v-else>
        Not logged in
      </div>
    </div>

    <button
      class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
      @click="spotifyRedirect"
    >
      Login
    </button>
  </div>
</template>