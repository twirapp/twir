import { useLocalStorage, useViewportSize } from '@mantine/hooks';
import React, { useEffect, useRef, useState } from 'react';
import Moveable from 'react-moveable';

import { BotManage } from './bot-manage';

const MoveableElement: React.FC<React.PropsWithChildren<{
  viewportHeight: number,
  bound: { left: number, right: number, top: number, bottom: number }
  widgetName: string,
}>> = (props) => {
  const targetRef = React.useRef<HTMLDivElement>(null);
  const [value, setValue] = useLocalStorage({
    key: `${props.widgetName}-widget-pos`,
    defaultValue: {
      height: 250,
      width: 400,
      transform: '',
    },
  });

  return (
    <div>
      <div ref={targetRef} style={{ height: value.height, width: value.width, transform: value.transform }}>
        {props.children}
      </div>

      <Moveable
        target={targetRef}
        draggable={true}
        resizable={true}
        snappable={true}
        bounds={{ ...props.bound }}
        onResize={(e) => {
          const beforeTranslate = e.drag.beforeTranslate;
          e.target.style.width = `${e.width}px`;
          e.target.style.height = `${e.height}px`;
          e.target.style.transform = `translate(${beforeTranslate[0]}px, ${beforeTranslate[1]}px)`;
          setValue({ width: e.width, height: e.height, transform: e.target.style.transform });
        }}
        onDrag={e => {
          e.target.style.transform = `translate(${e.beforeTranslate[0]}px, ${e.beforeTranslate[1]}px)`;
          setValue({ width: e.width, height: e.height, transform: e.target.style.transform });
        }}
      />
    </div>
  );
};


export const DashboardWidgets = () => {
  const dragAreaRef = useRef<HTMLDivElement>(null);
  const { height: viewportHeight } = useViewportSize();
  const [bound, setBound] = useState<{ left: number, right: number, top: number, bottom: number }>({
    left: 0,
    right: 0,
    bottom: 0,
    top: 0,
  });

  useEffect(() => {
    if (!dragAreaRef.current) return;
    const {
      left,
      right,
      top,
      bottom,
    } = dragAreaRef.current.getBoundingClientRect();
    setBound({ left, right, top, bottom });
  }, [dragAreaRef.current]);

  return (
    <div ref={dragAreaRef} style={{ height: viewportHeight }}>
     <MoveableElement viewportHeight={viewportHeight} bound={bound} widgetName={'bot-manage'}>
       <BotManage />
     </MoveableElement>

      <MoveableElement viewportHeight={viewportHeight} bound={bound} widgetName={'bot-manage2'}>
        <BotManage />
      </MoveableElement>

    </div>
  );
};
