import { css } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';
import { CanvasConnection } from 'app/features/canvas/element';
import { ElementState } from 'app/features/canvas/runtime/element';
import { Scene } from 'app/features/canvas/runtime/scene';

type Props = {
  setSVGRef: (anchorElement: SVGSVGElement) => void;
  setLineRef: (anchorElement: SVGLineElement) => void;
  scene: Scene;
};

interface ConnectionInfo {
  source: ElementState;
  target: ElementState;
  info: CanvasConnection;
}

export const ConnectionSVG = ({ setSVGRef, setLineRef, scene }: Props) => {
  const styles = useStyles2(getStyles);

  let selectedConnection: CanvasConnection | undefined = undefined;
  const onKeyUp = (e: KeyboardEvent) => {
    if (e.keyCode === 8) {
      if (selectedConnection) {
        const sourceElement = scene.byName.get(selectedConnection?.sourceName ?? '');

        if (sourceElement) {
          sourceElement.options.connections = sourceElement!.options.connections?.filter(
            (connection) => connection !== selectedConnection
          );
          sourceElement.onChange(sourceElement.options);
        }

        selectedConnection = undefined;
      }
    }

    window.document.removeEventListener('keyup', onKeyUp);
    scene.selecto!.rootContainer!.removeEventListener('click', clearSelectedConnection);
  };

  const clearSelectedConnection = (event: MouseEvent) => {
    // TODO: Handle case where clicking on connections
    selectedConnection = undefined;
  };

  const selectConnection = (connection: CanvasConnection) => {
    selectedConnection = connection;

    window.document.addEventListener('keyup', onKeyUp);
    scene.selecto!.rootContainer!.addEventListener('click', clearSelectedConnection);
  };

  // TODO: memos? in scene?  only update when things actually change?
  // Flat list of all connections
  const findConnections = () => {
    const connections: ConnectionInfo[] = [];
    for (let v of scene.byName.values()) {
      if (v.options.connections) {
        for (let c of v.options.connections) {
          const target = c.targetName ? scene.byName.get(c.targetName) : v.parent;
          if (target) {
            connections.push({
              source: v,
              target,
              info: c,
            });
          }
        }
      }
    }
    return connections;
  };

  // Figure out target and then target's relative coordinates drawing (if no target do parent)
  const renderConnections = () => {
    return findConnections().map((v, idx) => {
      const { source, target, info } = v;
      const sourceRect = source.div?.getBoundingClientRect();
      const parent = source.div?.parentElement;
      const parentRect = parent?.getBoundingClientRect();

      if (!sourceRect || !parent || !parentRect) {
        return;
      }

      const parentBorderWidth = parseFloat(getComputedStyle(parent).borderWidth);

      const sourceHorizontalCenter = sourceRect.left - parentRect.left - parentBorderWidth + sourceRect.width / 2;
      const sourceVerticalCenter = sourceRect.top - parentRect.top - parentBorderWidth + sourceRect.height / 2;

      const x1 = sourceHorizontalCenter + (info.source.x * sourceRect.width) / 2;
      const y1 = sourceVerticalCenter - (info.source.y * sourceRect.height) / 2;

      let x2;
      let y2;

      if (info.targetName) {
        const targetRect = target.div?.getBoundingClientRect();

        const targetHorizontalCenter = targetRect!.left - parentRect.left - parentBorderWidth + targetRect!.width / 2;
        const targetVerticalCenter = targetRect!.top - parentRect.top - parentBorderWidth + targetRect!.height / 2;

        x2 = targetHorizontalCenter + (info.target.x * targetRect!.width) / 2;
        y2 = targetVerticalCenter - (info.target.y * targetRect!.height) / 2;
      } else {
        const parentHorizontalCenter = parentRect.width / 2;
        const parentVerticalCenter = parentRect.height / 2;

        x2 = parentHorizontalCenter + (info.target.x * parentRect.width) / 2;
        y2 = parentVerticalCenter - (info.target.y * parentRect.height) / 2;
      }

      // TODO: Handle showing visual feedback when selected
      return (
        <svg className={styles.connection} key={idx}>
          <g onClick={() => selectConnection(info)}>
            <defs>
              <marker
                id="head"
                markerWidth="10"
                markerHeight="7"
                refX="10"
                refY="3.5"
                orient="auto"
                stroke="rgb(255,255,255)"
                pointerEvents="auto"
              >
                <polygon points="0 0, 10 3.5, 0 7" fill="rgb(255,255,255)" />
              </marker>
            </defs>
            <line
              stroke="rgb(255,255,255)"
              pointerEvents="auto"
              strokeWidth={2}
              markerEnd="url(#head)"
              x1={x1}
              y1={y1}
              x2={x2}
              y2={y2}
            />
          </g>
        </svg>
      );
    });
  };

  return (
    <>
      <svg ref={setSVGRef} className={styles.editorSVG}>
        <defs>
          <marker
            id="editorHead"
            markerWidth="10"
            markerHeight="7"
            refX="10"
            refY="3.5"
            orient="auto"
            stroke="rgb(255,255,255)"
          >
            <polygon points="0 0, 10 3.5, 0 7" fill="rgb(255,255,255)" />
          </marker>
        </defs>
        <line ref={setLineRef} stroke="rgb(255,255,255)" strokeWidth={2} markerEnd="url(#editorHead)" />
      </svg>
      {renderConnections()}
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  editorSVG: css`
    position: absolute;
    pointer-events: none;
    width: 100%;
    height: 100%;
    z-index: 1000;
    display: none;
  `,
  connection: css`
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 1000;
    pointer-events: none;
  `,
});
