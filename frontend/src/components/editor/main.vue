<template>
  <div class="editor">
    <div class="toolbox">
      <div
        class="element-template node"
        draggable="true"
        v-on:dragstart="onDragStart($event, 'node')"
      >
        Node
      </div>
      <div
        class="element-template resistor"
        draggable="true"
        v-on:dragstart="onDragStart($event, 'resistor')"
      >
        Resistor
      </div>
    </div>
    <canvas
      ref="canvas"
      v-on:dragover.prevent
      v-on:drop="onDrop"
      v-on:mousedown="onMouseDown"
      v-on:mousemove="onMouseMove"
      v-on:mouseup="onMouseUp"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Element, Connection } from './models'

export default defineComponent({
  name: 'CircuitEditor',
  data() {
    return {
      elements: [] as Element[],
      connections: [] as Connection[],
      nextId: 1,
      canvasSize: { width: 0, height: 0 },
      ctx: null as CanvasRenderingContext2D | null,
      selectedElement: null as Element | null,
      dragOffset: { x: 0, y: 0 },
      isDragging: false,
      isConnecting: false,
      connectingPort: null as string | null,
      hoveredPort: null as string | null,
      isRotating: false,
      rotationStartAngle: 0,
      initialAngle: 0,
    }
  },
  mounted() {
    this.initEditor()
    window.addEventListener('resize', this.initEditor)
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.initEditor)
  },

  computed: {
    canvas(): HTMLCanvasElement {
      return this.$refs.canvas as HTMLCanvasElement
    },
  },

  methods: {
    initEditor() {
      const container = this.canvas.parentElement
      if (container) {
        this.canvasSize.width = container.clientWidth
        this.canvasSize.height = container.clientHeight
        this.canvas.width = this.canvasSize.width
        this.canvas.height = this.canvasSize.height
      }
      this.ctx = this.canvas.getContext('2d')
      this.draw()
    },

    generateId() {
      return `el-${this.nextId++}`
    },

    onDragStart(event: DragEvent, type: 'node' | 'resistor') {
      // начинаем двигать элемент с левой панели
      console.log('on_drag_start', type)
      event.dataTransfer?.setData('text/plain', type)
    },

    onDrop(event: DragEvent) {
      // отпускаем элемент на canvas, для создания
      console.log('on_drop')
      if (!this.ctx) return

      const type = event.dataTransfer?.getData('text/plain') as
        | 'node'
        | 'resistor'
      const rect = (event.target as HTMLCanvasElement).getBoundingClientRect()
      const x = event.clientX - rect.left
      const y = event.clientY - rect.top

      const newElement: Element = {
        id: this.generateId(),
        type,
        x,
        y,
      }

      if (type === 'resistor') {
        newElement.ports = [
          { x: -20, y: 0 },
          { x: 20, y: 0 },
        ]
        newElement.connections = []
      } else {
        newElement.connections = []
      }

      this.elements.push(newElement)
      this.draw()
    },

    onMouseDown(event: MouseEvent) {
      // Зажимаем курсор
      if (!this.ctx) return
      const rect = this.canvas.getBoundingClientRect()
      // получим текущее положение курсора
      const mouseX = event.clientX - rect.left
      const mouseY = event.clientY - rect.top
      // Снимаем выделение со всех элементов
      this.elements.forEach((el) => (el.selected = false))
      // Провелим, мы зажали на порт или нет
      const portInfo = this.findPortAtPosition(mouseX, mouseY)
      if (portInfo) {
        this.isConnecting = true
        this.connectingPort = `${portInfo.element.id}:${portInfo.portIndex}`
        return
      }
      // Если зажали на элемент, начинаем отслеживать его перемещение
      const element = this.findElementAtPosition(mouseX, mouseY)
      if (element) {
        // Выделяем элемент
        element.selected = true
        this.selectedElement = element
        this.isDragging = true
        this.dragOffset = {
          x: mouseX - element.x,
          y: mouseY - element.y,
        }
        return
      }
      // Если кликнули на пустое место - снимаем выделение
      this.selectedElement = null
      this.draw()
    },

    drawOrthogonalConnection(
      ctx: CanvasRenderingContext2D,
      startX: number,
      startY: number,
      endX: number,
      endY: number,
    ) {
      // идет пересчет линки соединительной
      // Определяем промежуточные точки для создания соединения с прямыми углами
      const midX = startX + (endX - startX) / 2
      const midY = startY + (endY - startY) / 2
      // Сначала горизонтально, затем вертикально
      ctx.beginPath()
      ctx.moveTo(startX, startY)
      ctx.lineTo(midX, startY) // горизонтальная линия
      ctx.lineTo(midX, endY) // вертикальная линия
      ctx.lineTo(endX, endY) // горизонтальная линия
      ctx.stroke()
    },

    onMouseMove(event: MouseEvent) {
      // Перемещение курсора по канвасу
      if (!this.ctx) return
      const rect = this.canvas.getBoundingClientRect()
      const mouseX = event.clientX - rect.left
      const mouseY = event.clientY - rect.top
      // Update hovered port for visual feedback
      const portInfo = this.findPortAtPosition(mouseX, mouseY)
      this.hoveredPort = portInfo
        ? `${portInfo.element.id}:${portInfo.portIndex}`
        : null
      if (this.isDragging && this.selectedElement) {
        this.selectedElement.x = mouseX - this.dragOffset.x
        this.selectedElement.y = mouseY - this.dragOffset.y
        // Update connections positions
        this.draw()
      } else if (this.isConnecting && this.connectingPort) {
        // Draw temporary connection line
        this.draw()
        this.drawTemporaryConnection(mouseX, mouseY)
      }
    },
    removeElement(elementId: string) {
      // Удаляем все соединения, связанные с этим элементом
      this.connections = this.connections.filter(
        (conn) => !conn.from.startsWith(elementId) && conn.to !== elementId,
      )

      // Удаляем сам элемент
      this.elements = this.elements.filter((el) => el.id !== elementId)

      this.draw()
    },

    onMouseUp(event: MouseEvent) {
      this.isRotating = false
      if (!this.ctx) return

      const canvas = this.$refs.canvas as HTMLCanvasElement
      const rect = canvas.getBoundingClientRect()
      const mouseX = event.clientX - rect.left
      const mouseY = event.clientY - rect.top

      if (this.isConnecting && this.connectingPort) {
        // Проверяем, что соединяем с узлом
        const targetElement = this.findElementAtPosition(mouseX, mouseY)

        if (targetElement && targetElement.type === 'node') {
          const [sourceElementId, portIndex] = this.connectingPort.split(':')

          // Проверяем, не соединен ли уже этот порт
          const existingConnectionIndex = this.connections.findIndex(
            (conn) => conn.from === this.connectingPort,
          )

          // Если порт уже соединен, сначала удаляем старое соединение
          if (existingConnectionIndex !== -1) {
            this.connections.splice(existingConnectionIndex, 1)

            // Также удаляем ссылку в элементе
            const sourceElement = this.elements.find(
              (el) => el.id === sourceElementId,
            )
            if (sourceElement && sourceElement.connections) {
              //@ts-ignore
              sourceElement.connections[parseInt(portIndex)] = undefined
            }
          }

          // Создаем новое соединение
          const newConnection: Connection = {
            id: this.generateId(),
            from: this.connectingPort,
            to: targetElement.id,
          }
          this.connections.push(newConnection)

          // Обновляем связи в элементах
          const sourceElement = this.elements.find(
            (el) => el.id === sourceElementId,
          )
          if (sourceElement && sourceElement.connections) {
            sourceElement.connections[parseInt(portIndex)] = targetElement.id
          }
        }
      }

      this.isDragging = false
      this.isConnecting = false
      this.selectedElement = null
      this.connectingPort = null
      this.draw()
    },

    findElementAtPosition(x: number, y: number): Element | null {
      // Check in reverse order to prioritize elements on top
      for (let i = this.elements.length - 1; i >= 0; i--) {
        const element = this.elements[i]
        const width = element.type === 'node' ? 20 : 40
        const height = element.type === 'node' ? 20 : 20

        if (
          x >= element.x - width / 2 &&
          x <= element.x + width / 2 &&
          y >= element.y - height / 2 &&
          y <= element.y + height / 2
        ) {
          return element
        }
      }
      return null
    },

    findPortAtPosition(
      x: number,
      y: number,
    ): { element: Element; portIndex: number } | null {
      for (const element of this.elements) {
        if (element.ports) {
          for (let i = 0; i < element.ports.length; i++) {
            const port = element.ports[i]
            const portX = element.x + port.x
            const portY = element.y + port.y

            // Check if mouse is within port circle (radius 5)
            if (Math.sqrt((x - portX) ** 2 + (y - portY) ** 2) <= 5) {
              return { element, portIndex: i }
            }
          }
        }
      }
      return null
    },

    draw() {
      if (!this.ctx) return
      console.log('С нуля отрисовываем канвас?')
      const ctx = this.ctx
      ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)

      // Draw connections first (so they appear under elements)
      this.drawConnections()

      // Draw elements
      this.elements.forEach((element) => {
        if (element.type === 'node') {
          this.drawNode(element)
        } else if (element.type === 'resistor') {
          this.drawResistor(element)
        }
      })
    },

    drawConnections() {
      if (!this.ctx) return
      const ctx = this.ctx

      ctx.strokeStyle = '#000'
      ctx.lineWidth = 2

      this.connections.forEach((connection) => {
        const [fromElementId, portIndex] = connection.from.split(':')
        const fromElement = this.elements.find((el) => el.id === fromElementId)
        const toElement = this.elements.find((el) => el.id === connection.to)

        if (fromElement && toElement && fromElement.ports) {
          const port = fromElement.ports[parseInt(portIndex)]
          const startX = fromElement.x + port.x
          const startY = fromElement.y + port.y

          this.drawOrthogonalConnection(
            ctx,
            startX,
            startY,
            toElement.x,
            toElement.y,
          )
        }
      })
    },

    drawTemporaryConnection(x: number, y: number) {
      if (!this.ctx || !this.connectingPort) return
      const ctx = this.ctx

      const [elementId, portIndex] = this.connectingPort.split(':')
      const fromElement = this.elements.find((el) => el.id === elementId)

      if (fromElement && fromElement.ports) {
        const port = fromElement.ports[parseInt(portIndex)]
        const startX = fromElement.x + port.x
        const startY = fromElement.y + port.y

        ctx.strokeStyle = '#888'
        ctx.lineWidth = 2
        ctx.setLineDash([5, 5])
        this.drawOrthogonalConnection(ctx, startX, startY, x, y)
        ctx.setLineDash([])
      }
    },

    drawResistor(resistor: Element) {
      if (!this.ctx || !resistor.ports) return
      const ctx = this.ctx

      // Рисуем обводку если элемент выделен
      if (resistor.selected) {
        ctx.strokeStyle = '#FF5722'
        ctx.lineWidth = 3
        ctx.strokeRect(resistor.x - 25, resistor.y - 15, 50, 30)
      }

      // Остальная отрисовка резистора
      ctx.fillStyle = '#FF9800'
      ctx.fillRect(resistor.x - 20, resistor.y - 10, 40, 20)
      ctx.strokeStyle = '#000'
      ctx.lineWidth = 1
      ctx.strokeRect(resistor.x - 20, resistor.y - 10, 40, 20)

      // Отрисовка портов
      resistor.ports.forEach((port, index) => {
        const portX = resistor.x + port.x
        const portY = resistor.y + port.y

        const portId = `${resistor.id}:${index}`
        if (portId === this.hoveredPort || portId === this.connectingPort) {
          ctx.fillStyle = '#FF5722'
        } else {
          ctx.fillStyle = '#2196F3'
        }

        ctx.beginPath()
        ctx.arc(portX, portY, 5, 0, Math.PI * 2)
        ctx.fill()
        ctx.stroke()
      })
    },

    drawNode(node: Element) {
      if (!this.ctx) return
      const ctx = this.ctx

      // Рисуем обводку если элемент выделен
      if (node.selected) {
        ctx.strokeStyle = '#FF5722'
        ctx.lineWidth = 3
        ctx.beginPath()
        ctx.arc(node.x, node.y, 15, 0, Math.PI * 2)
        ctx.stroke()
      }

      // Остальная отрисовка узла
      ctx.fillStyle = '#4CAF50'
      ctx.beginPath()
      ctx.arc(node.x, node.y, 10, 0, Math.PI * 2)
      ctx.fill()
      ctx.strokeStyle = '#000'
      ctx.lineWidth = 1
      ctx.stroke()

      if (
        this.isConnecting &&
        this.findElementAtPosition(node.x, node.y)?.id === node.id
      ) {
        ctx.strokeStyle = '#FF5722'
        ctx.lineWidth = 2
        ctx.beginPath()
        ctx.arc(node.x, node.y, 12, 0, Math.PI * 2)
        ctx.stroke()
      }
    },
  },
})
</script>

<style lang="scss">
.editor {
  height: 100%;
  width: 100%;
  background-color: antiquewhite;
  display: flex;

  .toolbox {
    width: 150px;
    background-color: #f0f0f0;
    padding: 10px;
    border-right: 1px solid #ccc;

    .element-template {
      padding: 10px;
      margin-bottom: 10px;
      background-color: white;
      border: 1px solid #999;
      border-radius: 4px;
      cursor: grab;
      text-align: center;

      &.node {
        background-color: #4caf50;
        color: white;
      }

      &.resistor {
        background-color: #ff9800;
        color: white;
      }

      &:hover {
        opacity: 0.8;
      }
    }
  }

  > canvas {
    flex-grow: 1;
    height: 100%;
    background-color: white;
    background-size: 20px 20px;
    background-image: linear-gradient(
        to right,
        #e0e0e0 0.5px,
        transparent 0.5px
      ),
      linear-gradient(to bottom, #e0e0e0 0.5px, transparent 0.5px);
  }
}
</style>
