<template>
  <div class="new_task">
    <div class="btn_new_calc">
      <div class="selector_data">
        <div class="btn_item" :class="{ is_clicked: clicked_btn == 'nodes' }" v-on:click="handler_click_btn('nodes')">
          Узлы
        </div>
        <div class="btn_item" :class="{ is_clicked: clicked_btn == 'branches' }"
          v-on:click="handler_click_btn('branches')">
          Ветви
        </div>
        <div class="btn_item" :class="{ is_clicked: clicked_btn == 'results', disabled: clicked_btn != 'results' }"
          v-on:click="handler_click_btn('results')">
          Результаты
        </div>
      </div>
      <div class="selector_data task">
        <div class="btn_item" v-on:click="handler_click_btn_start">
          Запустить расчет
        </div>
      </div>
    </div>
    <div class="_tbl_container" v-if="clicked_btn == 'nodes'">
      <div class="tbl">
        <div class="tbl_header tbl_row">
          <div v-for="(item, index) in nodes.labels" :key="index" class="tbl_row_itm" :class="get_class_row(item, index)">
            {{ (translates as any)[item] }}
          </div>
        </div>
        <div class="tbl_content">
          <div v-for="(item, key) in nodes.items" class="tbl_row">
            <div v-for="(label, index) in nodes.labels" :key="key" class="tbl_row_itm"
              :class="get_class_row(label, index)">
              <input type=" text" v-model="(item as any)[label]" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="_tbl_container" v-else-if="clicked_btn == 'branches'">
      <div class="tbl">
        <div class="tbl_header tbl_row">
          <div v-for="(item, index) in branches.labels" :key="index" class="tbl_row_itm">
            {{ (translates as any)[item] }}
          </div>
        </div>
        <div class="tbl_content">
          <div v-for="(item, key) in branches.items" class="tbl_row">
            <div v-for="(label, index) in branches.labels" :key="key" class="tbl_row_itm" :class="{ first: index == 0 }">
              <input type="text" v-model="(item as any)[label]" />
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="_tbl_container" v-else-if="clicked_btn == 'results'">
      <div class="selector_data">
        <div class="btn_item">
          Напряжения в узлах
        </div>
      </div>
      <div class="tbl disabled">
        <div class="tbl_header tbl_row">
          <div v-for="(item, index) in result.node_labels" :key="index" class="tbl_row_itm"
            :class="get_class_row(item, index)">
            {{ (translates as any)[item] }}
          </div>
        </div>
        <div class="tbl_content">
          <div v-for="(item, key) in nodes.items" class="tbl_row">
            <div v-for="(label, index) in result.node_labels" :key="key" class="tbl_row_itm"
              :class="get_class_row(label, index)">
              <input v-if="label == 'voltage'" type="text" disabled
                v-model="(result.calculated as any).Nodes[key].Voltage" />
              <input v-else type="text" disabled v-model="(item as any)[label]" />
            </div>
          </div>
        </div>
      </div>
      <div class="selector_data">
        <div class="btn_item">
          Токи в ветвях
        </div>
      </div>
      <div class="tbl disabled">
        <div class="tbl_header tbl_row">
          <div v-for="(item, index) in result.branch_labels" :key="index" class="tbl_row_itm">
            {{ (translates as any)[item] }}
          </div>
        </div>
        <div class="tbl_content">
          <div v-for="(item, key) in branches.items" class="tbl_row">
            <div v-for="(label, index) in result.branch_labels" :key="key" class="tbl_row_itm"
              :class="{ first: index == 0 }">
              <input v-if="label == 'current'" type="text" disabled
                v-model="(result.calculated as any).Branches[key].Current" />
              <input v-else type="text" disabled v-model="(item as any)[label]" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Calculate } from '../../wailsjs/go/main/App'

export default defineComponent({
  components: {},
  data() {
    return {
      result: {
        node_labels: ['number', 'name', 'voltage'],
        branch_labels: ['number', 'name', 'node_first', 'node_second', 'current'],
        calculated: {}
      },
      translates: {
        number: '№',
        name: 'Имя',
        nominalvoltage: 'Uном, кВ',
        voltage: 'U, кВ',
        current: 'I, А',
        node_first: 'Начало, №',
        node_second: 'Конец, №',
        r: 'R, Ом',
        x: 'X, Ом',
        activepower: 'P, МВт',
        reactivepower: 'Q, МВАр'
      },
      nodes: {
        labels: ['number', 'name', 'nominalvoltage', 'activepower', 'reactivepower'],
        items: {
          '1': {
            number: '1',
            name: 'Узел1',
            nominalvoltage: '115',
            activepower: "0.1",
            reactivepower: '0.2'
          },
          '2': {
            number: '2',
            name: 'Узел2',
            nominalvoltage: '115',
            activepower: "0.0",
            reactivepower: '0.0'
          },
          '3': {
            number: '3',
            name: 'Узел3',
            nominalvoltage: '115',
            activepower: "0.9",
            reactivepower: '-1.0'
          },
        },
      },
      branches: {
        labels: ['number', 'name', 'node_first', 'node_second', 'r', 'x'],
        items: {
          '1': {
            number: '1',
            name: 'Линия1',
            node_first: '1',
            node_second: '2',
            r: '10',
            x: '20'
          },
          '2': {
            number: '2',
            name: 'Линия2',
            node_first: '2',
            node_second: '3',
            r: '5',
            x: '15'
          },
        }

      },
      clicked_btn: 'nodes'
    }
  },

  computed: {
    is_calculated() {
      return Object.keys(this.result.calculated).length != 0
    }
  },

  methods: {
    handler_click_btn(btn_type: string) {
      console.log('handler_click_btn', this.is_calculated)
      if (['branches', 'nodes'].includes(btn_type)) {
        this.result.calculated = {}
        this.clicked_btn = btn_type
      } else {
        if (this.is_calculated) {
          this.clicked_btn = btn_type
        } else {
          return
        }
      }
    },

    async handler_click_btn_start() {
      const res: any = await Calculate('start_task')
      this.result.calculated = res
      this.handler_click_btn('results')

    },

    get_class_row(row_name: string, index: number) {
      if (index == 0) {
        return row_name + ' first'
      }
      return row_name
    }
  }
})
</script>

<style lang="scss">
._tbl_container {
  width: 100%;
}

.new_task {
  height: 100%;
  width: 100%;

  >.btn_new_calc {
    display: flex;
    justify-content: space-between;


    >.selector_data {
      display: flex;
      width: 300px;
      background-color: #fafafa;
      border-radius: 5px;
      display: flex;
      align-items: center;
      margin-top: 20px;
      margin-left: 10px;
      overflow: hidden;

      >.btn_item {
        height: 30px;
        width: 150px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;

        &.is_clicked {
          background-color: #eae4e4;
        }

        &.disabled {
          cursor: not-allowed;
        }

      }

      &.task {
        margin-right: 10px;
        background-color: #808181;
        color: white;
        width: 150px;

        >.btn_item {
          width: 100%;
        }

      }
    }
  }
}

input {
  border: unset;
  background: unset;
  height: calc(100% - 5px);
  max-width: calc(100% - 10px);

  &:focus {
    outline: none;
    background-color: white;
  }
}

.tbl {
  padding: 10px;
  height: calc(100% - 80px);
  width: calc(100% - 20px);
  // overflow-x: auto;


  .tbl_row {
    display: flex;
    width: 100%;
    flex: 0 0 auto;
    border-bottom: 1px solid #eae4e4;

    ;

    >.tbl_row_itm {
      min-width: 80px;
      flex-grow: 1;
      flex-basis: 0;
      display: flex;
      padding-left: 10px;
      height: 30px;
      align-items: center;
      border-left: 1px solid #eae4e4;

      &.first {
        border-left: unset;
      }

      &.number {
        width: 50px !important;
      }
    }

    &.tbl_header {
      background-color: #e4e7e7;
      border-radius: 5px 5px 0 0;
      overflow: hidden;


      >.tbl_row_itm {
        border-right: 1px solid #eae4e4;
      }
    }
  }

  >.tbl_content {
    border: 1px solid #eae4e4;
    height: 100%;
    min-height: 130px;
    background-color: #fafafa;
    border-radius: 0 0 5px 5px;
    width: calc(100% - 2px);
  }

  &.disabled {
    input {
      cursor: not-allowed;
    }
  }
}
</style>
