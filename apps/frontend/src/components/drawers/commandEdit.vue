<script lang="ts" setup>
import { mdiPlus, mdiClose, mdiApplicationVariable } from '@mdi/js';
import { useStore } from '@nanostores/vue';
import { ref } from 'vue';
import { useDisplay } from 'vuetify';

import confirmDeletion from '@/components/confirmDeletion.vue';
import { editableCommand } from '@/stores/commands';
import { variablesStore } from '@/stores/variables';

const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const addAliaseInput = ref('');
const variables = useStore(variablesStore);
const responsesRef = ref<HTMLTextAreaElement[]>([]);
const { smAndDown } = useDisplay();

const symbolsRegexp = /^\W$|_|/;
function preventSymbolsInCommandName(e: KeyboardEvent) {
  // if (symbolsRegexp.test(e.key)) {
  //   e.preventDefault();
  // }
}

function addAliase() {
  editableCommand.value!.aliases.push(addAliaseInput.value);
  addAliaseInput.value = '';
}

function onDelete() {
  console.log(editableCommand.value);
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
          v-model="editableCommand!.name" 
          prefix="!" 
          label="Name" 
          :rules="[
            v => !!v || 'Field is required'
          ]"
          @keydown="preventSymbolsInCommandName" 
        />
        <v-text-field 
          v-model="addAliaseInput" 
          prefix="!"
          label="New aliase"
          :append-icon="mdiPlus"
          :rules="[
            (v) => (!editableCommand!.aliases.includes(v) && editableCommand!.name != v) || 'Aliase already exists'
          ]"
          @click:append="addAliase"
          @keydown="preventSymbolsInCommandName" 
          @keyup.enter="addAliase"
        />
        <div class="d-flex flex-wrap">
          <v-chip
            v-for="(aliase, index) of editableCommand!.aliases"
            :key="aliase"
            closable
            class="mr-2 mt-2"
            @click:close="editableCommand!.aliases.splice(index, 1)"
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
              @click="editableCommand!.responses!.push({
                order: editableCommand!.responses?.length ? editableCommand!.responses?.length - 1 : 0,
                text: '',
              } as any)"
            >
              <v-icon>{{ mdiPlus }}</v-icon>
            </v-btn>
          </div>
        </div>

        <v-textarea
          v-for="(response, responseIndex) of editableCommand!.responses" 
          :key="responseIndex"
          ref="responsesRef"
          v-model="editableCommand!.responses![responseIndex].text"
          auto-grow
          rows="1"
          row-height="5"
          class="mt-2"
          :rules="[
            v => !!v || 'Field is required'
          ]"
        >
          <!-- :append-icon="mdiClose" -->
          <!-- @click:append="command.responses!.splice(index, 1)" -->
          <template #append-inner>
            <v-menu max-height="200">
              <template #activator="{ props }">
                <v-btn
                  v-bind="props"
                  :icon="mdiApplicationVariable" variant="tonal" size="small" style="margin-top: -8px;"
                />
              </template>
              <v-list>
                <v-list-item
                  v-for="(variable, variableIndex) in variables"
                  :key="variableIndex"
                  :value="variableIndex"
                  @click="response.text += ` $(${variable.example ?? variable.name})`"
                >
                  <b>{{ `$(${variable.name})` }}</b>
                  <v-list-item-title>{{ variable.description }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </template>
          <template #append>
            <v-btn 
              style="margin-top: -8px" 
              variant="tonal" 
              :icon="mdiClose"
              size="small" 
              @click="editableCommand!.responses!.splice(responseIndex, 1)"
            />
          </template>
        </v-textarea>

        <div>
          <h4>Cooldown</h4>
          <div class="d-flex justify-space-between mt-2">
            <v-text-field 
              v-model="editableCommand!.cooldown"
              style="width: 45%"
              type="number"
              label="Seconds"
              :rules="[
                v => v >= 0 || 'Cannot be lower then 0'
              ]"
            />
            <v-select
              v-model="editableCommand!.cooldownType"
              style="width: 45%"
              class="ml-4"
              label="Cooldown Type"
              :items="['GLOBAL', 'USER']"
            ></v-select>
          </div>
        </div>
    
        <div>
          <v-select
            v-model="editableCommand!.permission"
            label="Permission"
            :items="['BROADCASTER', 'MODERATOR', 'SUBSCRIBER', 'VIP', 'VIEWER', 'FOLLOWER']"
          ></v-select>
        </div>

        <v-row>
          <v-col :cols="smAndDown ? 12 : 6">
            <v-checkbox 
              v-model="editableCommand!.visible" 
              label="Show command in list of commands" 
              density="compact"
            />
          </v-col>
          <v-col :cols="smAndDown ? 12 : 6">
            <v-checkbox 
              v-model="editableCommand!.keepResponsesOrder" 
              label="Keep order of commands responses" 
              density="compact"
            />
          </v-col>
          <v-col :cols="smAndDown ? 12 : 6">
            <v-checkbox 
              v-model="editableCommand!.isReply" 
              label="Use twitch reply feature" 
              density="compact"
            />
          </v-col>
        </v-row>

        <div class="mt-4">
          <v-textarea
            v-model="editableCommand!.description"
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

<style>
.v-checkbox {
  height: 20px;
}
</style>