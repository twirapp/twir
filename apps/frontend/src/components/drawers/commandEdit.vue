<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  command: ChannelCommand
}>();

const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const command = ref(props.command);
const addAliaseInput = ref('');

const symbolsRegexp = /^\W$|_|/;
function preventSymbolsInCommandName(e: KeyboardEvent) {
  // if (symbolsRegexp.test(e.key)) {
  //   e.preventDefault();
  // }
}

function addAliase() {
  command.value!.aliases.push(addAliaseInput.value);
  addAliaseInput.value = '';
}

function onDelete() {
  console.log(command.value);
  emits('cancel');
}
</script>

<template>
  <v-list-item>
    <div class="d-flex justify-space-between">
      <h1>Edit command</h1>
      <div class="d-flex d-inline">
        <v-btn size="small" class="mt-1 mr-2" @click="$emit('cancel')">
          Cancel
        </v-btn>
        <confirmDeletion :cb="() => onDelete()">
          <v-btn color="red" size="small" class="mt-1 mr-2">
            Delete
          </v-btn>
        </confirmDeletion>
        <v-btn color="green" size="small" class="mt-1">
          Save
        </v-btn>
      </div>
    </div>
  </v-list-item>

  <v-divider></v-divider>

  <v-list-item>
    <v-form class="mt-2">
      <div
        class="d-flex flex-column"
      >
        <v-text-field 
          v-model="command.name" 
          prefix="!" label="Name" 
          :rules="[
            v => !!v || 'Field is required'
          ]"
          @keydown="preventSymbolsInCommandName" 
        />
        <v-text-field 
          v-model="addAliaseInput" 
          prefix="!"
          label="Add aliase"
          :append-icon="mdiPlus"
          :rules="[
            (v) => (!command.aliases.includes(v) && command.name != v) || 'Aliase already exists'
          ]"
          @click:append="addAliase"
          @keydown="preventSymbolsInCommandName" 
          @keyup.enter="addAliase"
        />
        <div class="d-flex flex-wrap">
          <v-chip
            v-for="(aliase, index) of command!.aliases"
            :key="aliase"
            closable
            class="mr-2 mt-2"
            @click:close="command!.aliases.splice(index, 1)"
          >
            !{{ aliase }}
          </v-chip>
        </div>

        <v-divider class="mt-2"></v-divider>

        <div class="d-flex justify-space-between mt-2">
          <h4>Responses</h4>
          <div>
            <v-btn 
              variant="outlined" 
              size="x-small" 
              @click="command!.responses!.push({
                order: command!.responses?.length ? command!.responses?.length - 1 : 0,
                text: '',
              } as any)"
            >
              <v-icon>{{ mdiPlus }}</v-icon>
            </v-btn>
          </div>
        </div>

        <v-textarea
          v-for="(response, index) of command!.responses" 
          :key="index"
          v-model="command!.responses![index].text"
          auto-grow
          :append-icon="mdiClose"
          rows="1"
          row-height="5"
          class="mt-2"
          :rules="[
            v => !!v || 'Field is required'
          ]"
          @click:append="command.responses!.splice(index, 1)"
        />

        <div>
          <h4>Cooldown</h4>
          <div class="d-flex justify-space-between mt-2">
            <v-text-field 
              v-model="command.cooldown"
              style="width: 45%"
              type="number"
              label="Seconds"
              :rules="[
                v => v >= 0 || 'Cannot be lower then 0'
              ]"
            />
            <v-select
              v-model="command.cooldownType"
              style="width: 45%"
              class="ml-4"
              label="Cooldown Type"
              :items="['GLOBAL', 'USER']"
            ></v-select>
          </div>
        </div>
    
        <div>
          <v-select
            v-model="command.permission"
            label="Permission"
            :items="['BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER']"
          ></v-select>
        </div>

        <div>
          <v-checkbox v-model="command.visible" label="Show command in list of commands" />
          <v-checkbox v-model="command.keepResponsesOrder" label="Keep order of commands responses" />
          <v-checkbox v-model="command.isReply" label="Use twitch reply feature" />
        </div>

        <div>
          <v-textarea
            v-model="command.description"
            auto-grow
            label="Description"
            rows="1"
            row-height="5"
            class="mt-2"
          />
        </div>
      </div>
    </v-form>
  </v-list-item>
</template>