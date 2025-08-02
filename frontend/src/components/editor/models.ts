export interface Point {
  x: number
  y: number
}

export interface Element {
  id: string
  type: 'node' | 'resistor'
  x: number
  y: number
  ports?: Point[]
  angle?: number
  connections?: string[]
  selected?: boolean
}

export interface Connection {
  id: string
  from: string // elementId:portIndex
  to: string // elementId
}
