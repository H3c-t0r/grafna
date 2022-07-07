import { CSSProperties } from 'react';

import { ElementState } from 'app/features/canvas/runtime/element';
import { FrameState } from 'app/features/canvas/runtime/frame';
import { RootElement } from 'app/features/canvas/runtime/root';

export interface FlatElement {
  node: ElementState;
  depth: number;
  isOpen?: boolean;
}

export interface TreeElement {
  key: number;
  title: string;
  selectable?: boolean;
  children?: TreeElement[];
  dataRef: ElementState | FrameState;
  style?: CSSProperties;
}

export function getTreeData(root?: RootElement | FrameState, selection?: string[], selectedColor?: string) {
  let elements: TreeElement[] = [];
  if (root) {
    for (let i = root.elements.length; i--; i >= 0) {
      const item = root.elements[i];
      const element: TreeElement = {
        key: item.UID,
        title: item.getName(),
        selectable: true,
        dataRef: item,
      };

      const isSelected = isItemSelected(item, selection);
      if (isSelected) {
        element.style = { backgroundColor: selectedColor };
      }

      if (item instanceof FrameState) {
        element.children = getTreeData(item, selection, selectedColor);
      }
      elements.push(element);
    }
  }

  return elements;
}

function isItemSelected(item: ElementState, selection: string[] | undefined) {
  return Boolean(selection?.includes(item.getName()));
}
