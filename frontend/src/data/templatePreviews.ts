// SVG preview thumbnails for CV templates.
// Each preview is a stylized mini representation of the template layout.

type PreviewMap = Record<string, string>

function makeSvg(body: string): string {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 160 220" width="100%" height="100%">${body}</svg>`
}

const rect = (x: number, y: number, w: number, h: number, fill: string, rx = 0): string =>
  `<rect x="${x}" y="${y}" width="${w}" height="${h}" fill="${fill}"${rx ? ` rx="${rx}"` : ''}/>`

const circle = (cx: number, cy: number, r: number, fill: string): string =>
  `<circle cx="${cx}" cy="${cy}" r="${r}" fill="${fill}"/>`

// Reusable text lines
const line = (x: number, y: number, w: number, color = '#cbd5e1', h = 3): string =>
  rect(x, y, w, h, color, 1)

// Section: multiple short lines
function lines(x: number, y: number, width: number, count: number, color = '#cbd5e1', gap = 6): string {
  let out = ''
  for (let i = 0; i < count; i++) {
    const w = width * (0.7 + (i % 3) * 0.1)
    out += line(x, y + i * gap, w, color)
  }
  return out
}

// ---------- Templates ----------

// Modern: dark sidebar + main content
const modern = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 58, 220, '#2c3e50') +
  circle(29, 30, 16, '#ffffff') +
  line(8, 56, 42, '#ffffff', 4) +
  line(12, 66, 34, '#94a3b8') +
  line(8, 90, 42, '#3498db', 3) +
  lines(8, 100, 42, 4, '#94a3b8', 6) +
  line(8, 134, 42, '#3498db', 3) +
  lines(8, 144, 42, 4, '#94a3b8', 6) +
  line(66, 14, 80, '#1e293b', 5) +
  line(66, 24, 60, '#64748b') +
  line(66, 44, 90, '#3498db', 3) +
  lines(66, 54, 90, 3, '#cbd5e1', 7) +
  line(66, 90, 90, '#3498db', 3) +
  lines(66, 100, 90, 4, '#cbd5e1', 7) +
  line(66, 140, 90, '#3498db', 3) +
  lines(66, 150, 90, 3, '#cbd5e1', 7)
)

// Classic: centered header with underlines
const classic = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(20, 12, 120, 1, '#334155') +
  rect(20, 34, 120, 1, '#334155') +
  line(50, 18, 60, '#1e293b', 5) +
  line(56, 26, 48, '#64748b') +
  line(20, 46, 120, '#334155', 2) +
  lines(20, 52, 120, 4, '#94a3b8', 7) +
  line(20, 86, 120, '#334155', 2) +
  lines(20, 92, 120, 4, '#94a3b8', 7) +
  line(20, 126, 120, '#334155', 2) +
  lines(20, 132, 120, 4, '#94a3b8', 7) +
  line(20, 166, 120, '#334155', 2) +
  lines(20, 172, 120, 3, '#94a3b8', 7)
)

// Professional: corporate header bar
const professional = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 44, '#1a237e') +
  line(12, 14, 80, '#ffffff', 5) +
  line(12, 24, 60, '#b3c5ff') +
  line(12, 32, 40, '#b3c5ff') +
  line(12, 60, 80, '#1a237e', 3) +
  lines(12, 70, 136, 4, '#cbd5e1', 7) +
  line(12, 108, 80, '#1a237e', 3) +
  lines(12, 118, 136, 4, '#cbd5e1', 7) +
  line(12, 156, 80, '#1a237e', 3) +
  lines(12, 166, 136, 3, '#cbd5e1', 7)
)

// Creative: colorful shapes
const creative = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 70, '#E91E63') +
  circle(30, 35, 18, '#ffffff') +
  line(56, 26, 80, '#ffffff', 5) +
  line(56, 38, 60, '#ffc4d6') +
  line(56, 46, 40, '#ffc4d6') +
  circle(20, 100, 4, '#E91E63') +
  line(30, 98, 110, '#E91E63', 3) +
  lines(30, 106, 110, 3, '#cbd5e1', 6) +
  circle(20, 140, 4, '#9C27B0') +
  line(30, 138, 110, '#9C27B0', 3) +
  lines(30, 146, 110, 3, '#cbd5e1', 6) +
  circle(20, 180, 4, '#FF9800') +
  line(30, 178, 110, '#FF9800', 3) +
  lines(30, 186, 110, 2, '#cbd5e1', 6)
)

// Elegant: gold accents, serif feel
const elegant = makeSvg(
  rect(0, 0, 160, 220, '#fafaf7') +
  rect(16, 14, 128, 1, '#9C27B0') +
  rect(16, 40, 128, 1, '#9C27B0') +
  line(40, 20, 80, '#2d1b4e', 5) +
  line(48, 30, 64, '#9C27B0') +
  line(16, 52, 30, '#9C27B0', 3) +
  lines(16, 60, 128, 4, '#b8a7c9', 7) +
  line(16, 96, 30, '#9C27B0', 3) +
  lines(16, 104, 128, 4, '#b8a7c9', 7) +
  line(16, 140, 30, '#9C27B0', 3) +
  lines(16, 148, 128, 3, '#b8a7c9', 7)
)

// Academic: formal, dense text
const academic = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  line(20, 16, 120, '#1e293b', 5) +
  line(40, 26, 80, '#64748b') +
  line(20, 36, 120, '#64748b') +
  rect(16, 46, 128, 1, '#4CAF50') +
  line(16, 54, 50, '#4CAF50', 3) +
  lines(16, 62, 128, 3, '#94a3b8', 6) +
  line(16, 90, 50, '#4CAF50', 3) +
  lines(16, 98, 128, 4, '#94a3b8', 6) +
  line(16, 134, 50, '#4CAF50', 3) +
  lines(16, 142, 128, 4, '#94a3b8', 6) +
  line(16, 178, 50, '#4CAF50', 3) +
  lines(16, 186, 128, 3, '#94a3b8', 6)
)

// Minimal: lots of whitespace
const minimal = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  line(24, 30, 60, '#1e293b', 5) +
  line(24, 42, 40, '#94a3b8') +
  line(24, 68, 20, '#1e293b', 2) +
  lines(24, 76, 110, 3, '#cbd5e1', 7) +
  line(24, 116, 20, '#1e293b', 2) +
  lines(24, 124, 110, 3, '#cbd5e1', 7) +
  line(24, 164, 20, '#1e293b', 2) +
  lines(24, 172, 110, 2, '#cbd5e1', 7)
)

// Executive: bold top bar with gold accent
const executive = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 56, '#1e293b') +
  rect(0, 56, 160, 3, '#FF9800') +
  line(14, 18, 90, '#ffffff', 5) +
  line(14, 30, 60, '#FF9800') +
  line(14, 40, 80, '#cbd5e1') +
  line(14, 74, 70, '#FF9800', 3) +
  lines(14, 84, 132, 4, '#cbd5e1', 7) +
  line(14, 122, 70, '#FF9800', 3) +
  lines(14, 132, 132, 4, '#cbd5e1', 7) +
  line(14, 170, 70, '#FF9800', 3) +
  lines(14, 180, 132, 3, '#cbd5e1', 7)
)

// Tech: code-like styling, mono feel
const tech = makeSvg(
  rect(0, 0, 160, 220, '#f8fafc') +
  rect(0, 0, 160, 40, '#0f172a') +
  circle(10, 12, 2, '#ef4444') +
  circle(18, 12, 2, '#f59e0b') +
  circle(26, 12, 2, '#22c55e') +
  line(10, 22, 100, '#00BCD4', 4) +
  line(10, 30, 80, '#64748b') +
  line(10, 56, 60, '#00BCD4', 3) +
  lines(10, 64, 140, 3, '#cbd5e1', 6) +
  rect(10, 92, 30, 10, '#e0f7fa', 2) +
  rect(44, 92, 30, 10, '#e0f7fa', 2) +
  rect(78, 92, 30, 10, '#e0f7fa', 2) +
  rect(112, 92, 30, 10, '#e0f7fa', 2) +
  line(10, 116, 60, '#00BCD4', 3) +
  lines(10, 124, 140, 4, '#cbd5e1', 6) +
  line(10, 160, 60, '#00BCD4', 3) +
  lines(10, 168, 140, 3, '#cbd5e1', 6)
)

// Compact: dense single page
const compact = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 30, '#3F51B5') +
  line(8, 10, 80, '#ffffff', 4) +
  line(8, 20, 60, '#c5cae9') +
  line(8, 38, 50, '#3F51B5', 2) +
  lines(8, 44, 144, 2, '#cbd5e1', 5) +
  line(8, 66, 50, '#3F51B5', 2) +
  lines(8, 72, 144, 2, '#cbd5e1', 5) +
  line(8, 94, 50, '#3F51B5', 2) +
  lines(8, 100, 144, 2, '#cbd5e1', 5) +
  line(8, 122, 50, '#3F51B5', 2) +
  lines(8, 128, 144, 2, '#cbd5e1', 5) +
  line(8, 150, 50, '#3F51B5', 2) +
  lines(8, 156, 144, 2, '#cbd5e1', 5) +
  line(8, 178, 50, '#3F51B5', 2) +
  lines(8, 184, 144, 2, '#cbd5e1', 5)
)

// Bold: strong typography, big name
const bold = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(10, 14, 140, 10, '#F44336') +
  line(10, 30, 100, '#1e293b', 4) +
  rect(10, 50, 140, 2, '#F44336') +
  line(10, 60, 50, '#F44336', 4) +
  lines(10, 70, 140, 3, '#cbd5e1', 7) +
  line(10, 106, 50, '#F44336', 4) +
  lines(10, 116, 140, 4, '#cbd5e1', 7) +
  line(10, 156, 50, '#F44336', 4) +
  lines(10, 166, 140, 3, '#cbd5e1', 7)
)

// Clean: lots of whitespace with teal accents
const clean = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  circle(30, 30, 14, '#009688') +
  line(54, 22, 80, '#1e293b', 5) +
  line(54, 34, 60, '#94a3b8') +
  rect(14, 54, 132, 1, '#009688') +
  line(14, 62, 40, '#009688', 3) +
  lines(14, 70, 132, 3, '#cbd5e1', 7) +
  line(14, 104, 40, '#009688', 3) +
  lines(14, 112, 132, 3, '#cbd5e1', 7) +
  line(14, 146, 40, '#009688', 3) +
  lines(14, 154, 132, 3, '#cbd5e1', 7)
)

// Premium: gold everything
const premium = makeSvg(
  rect(0, 0, 160, 220, '#fffdf7') +
  rect(0, 0, 160, 60, '#2d1b0e') +
  rect(0, 60, 160, 2, '#c0982b') +
  circle(30, 30, 16, '#c0982b') +
  line(54, 22, 80, '#ffffff', 5) +
  line(54, 34, 60, '#c0982b') +
  line(54, 44, 50, '#d4b559') +
  line(14, 76, 50, '#c0982b', 3) +
  lines(14, 86, 132, 4, '#a8906a', 7) +
  line(14, 124, 50, '#c0982b', 3) +
  lines(14, 134, 132, 4, '#a8906a', 7) +
  line(14, 172, 50, '#c0982b', 3) +
  lines(14, 182, 132, 3, '#a8906a', 7)
)

// Formal: document-style
const formal = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(12, 12, 136, 1, '#455A64') +
  rect(12, 40, 136, 1, '#455A64') +
  line(50, 20, 60, '#1e293b', 5) +
  line(56, 30, 48, '#64748b') +
  line(12, 48, 40, '#455A64', 3) +
  lines(12, 56, 136, 3, '#94a3b8', 7) +
  line(12, 88, 40, '#455A64', 3) +
  lines(12, 96, 136, 4, '#94a3b8', 7) +
  line(12, 136, 40, '#455A64', 3) +
  lines(12, 144, 136, 4, '#94a3b8', 7) +
  line(12, 184, 40, '#455A64', 3) +
  lines(12, 192, 136, 2, '#94a3b8', 7)
)

// Timeline: vertical timeline with dots
const timeline = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  line(14, 14, 120, '#1e293b', 5) +
  line(14, 24, 80, '#94a3b8') +
  rect(30, 50, 2, 150, '#FF5722') +
  circle(31, 60, 4, '#FF5722') +
  line(42, 56, 100, '#FF5722', 3) +
  lines(42, 64, 100, 2, '#cbd5e1', 6) +
  circle(31, 100, 4, '#FF5722') +
  line(42, 96, 100, '#FF5722', 3) +
  lines(42, 104, 100, 2, '#cbd5e1', 6) +
  circle(31, 140, 4, '#FF5722') +
  line(42, 136, 100, '#FF5722', 3) +
  lines(42, 144, 100, 2, '#cbd5e1', 6) +
  circle(31, 180, 4, '#FF5722') +
  line(42, 176, 100, '#FF5722', 3) +
  lines(42, 184, 100, 2, '#cbd5e1', 6)
)

// Infographic: charts and visual elements
const infographic = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 50, '#8BC34A') +
  circle(26, 25, 14, '#ffffff') +
  line(50, 18, 80, '#ffffff', 5) +
  line(50, 30, 60, '#dcedc8') +
  line(12, 66, 50, '#8BC34A', 3) +
  rect(12, 76, 20, 10, '#8BC34A') +
  rect(36, 76, 40, 10, '#aed581') +
  rect(80, 76, 30, 10, '#c5e1a5') +
  line(12, 100, 50, '#8BC34A', 3) +
  circle(30, 118, 12, '#8BC34A') +
  circle(30, 118, 8, '#ffffff') +
  lines(50, 110, 90, 2, '#cbd5e1', 6) +
  line(12, 146, 50, '#8BC34A', 3) +
  rect(12, 156, 136, 4, '#e8f5e9', 2) +
  rect(12, 156, 100, 4, '#8BC34A', 2) +
  rect(12, 166, 136, 4, '#e8f5e9', 2) +
  rect(12, 166, 80, 4, '#aed581', 2) +
  rect(12, 176, 136, 4, '#e8f5e9', 2) +
  rect(12, 176, 120, 4, '#c5e1a5', 2)
)

// Additional templates that exist in CVFormView
const standard = classic
const ats = minimal
const teacher = academic
const medical = makeSvg(
  rect(0, 0, 160, 220, '#ffffff') +
  rect(0, 0, 160, 50, '#0077b6') +
  circle(80, 25, 10, '#ffffff') +
  rect(76, 20, 8, 2, '#0077b6') +
  rect(79, 17, 2, 8, '#0077b6') +
  line(12, 62, 50, '#0077b6', 3) +
  lines(12, 70, 136, 3, '#cbd5e1', 7) +
  line(12, 104, 50, '#0077b6', 3) +
  lines(12, 112, 136, 4, '#cbd5e1', 7) +
  line(12, 152, 50, '#0077b6', 3) +
  lines(12, 160, 136, 3, '#cbd5e1', 7)
)
const lawyer = makeSvg(
  rect(0, 0, 160, 220, '#fdfbf7') +
  rect(12, 12, 136, 1, '#78350f') +
  rect(12, 40, 136, 1, '#78350f') +
  line(44, 20, 72, '#2c1810', 5) +
  line(52, 30, 56, '#78350f') +
  line(12, 50, 40, '#78350f', 3) +
  lines(12, 58, 136, 3, '#94a3b8', 7) +
  line(12, 90, 40, '#78350f', 3) +
  lines(12, 98, 136, 4, '#94a3b8', 7) +
  line(12, 138, 40, '#78350f', 3) +
  lines(12, 146, 136, 4, '#94a3b8', 7) +
  line(12, 186, 40, '#78350f', 3) +
  lines(12, 194, 136, 2, '#94a3b8', 7)
)
const engineer = tech
const designer = creative
const minimalist = minimal

export const templatePreviews: PreviewMap = {
  modern,
  classic,
  professional,
  creative,
  elegant,
  academic,
  minimal,
  minimalist,
  executive,
  tech,
  engineer,
  compact,
  bold,
  clean,
  premium,
  formal,
  timeline,
  infographic,
  standard,
  ats,
  teacher,
  medical,
  lawyer,
  designer,
}

export function getTemplatePreview(id: string): string {
  return templatePreviews[id] || modern
}

export function getTemplatePreviewDataUrl(id: string): string {
  const svg = getTemplatePreview(id)
  return `data:image/svg+xml;utf8,${encodeURIComponent(svg)}`
}
