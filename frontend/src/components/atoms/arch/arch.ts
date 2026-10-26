/**
 * The mihrab arch, amami's signature shape: a pointed arch on straight jambs.
 * ARCH_PATH is drawn in a 40x50 box; ARCH_CLIP_PATH is the same outline in
 * objectBoundingBox units (0-1) for clipping images of any size.
 */
export const ARCH_PATH = 'M0.5 49.5 V22 C0.5 11 10 4 20 0.5 C30 4 39.5 11 39.5 22 V49.5 Z';

export const ARCH_CLIP_PATH = 'M0,1 V0.44 C0,0.22 0.24,0.08 0.5,0.01 C0.76,0.08 1,0.22 1,0.44 V1 Z';

export type ArchTone = 'sajadah' | 'bata' | 'kuningan' | 'outline';
