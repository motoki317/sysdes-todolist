<template>
  <q-page padding class='row'>

    <Suspense>
      <template #default>
        <todo-list></todo-list>
      </template>
      <template #fallback>
        <p>
          Loading...
        </p>
      </template>
    </Suspense>

  </q-page>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';
import TodoList from 'components/TodoList.vue';

export default defineComponent({
  name: 'PageIndex',

  components: { TodoList },

  setup() { return {} }
});
</script>
