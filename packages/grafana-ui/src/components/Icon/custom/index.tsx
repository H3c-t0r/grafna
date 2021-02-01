import React, { FC, ComponentType } from 'react';
import { SvgProps } from '../assets/types';

const InterpolationLinear: FC<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28.5 20" width={'30px'} height={size} {...rest}>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Icons">
          <circle cx="14.17" cy="2.67" r="2.67" />
          <circle cx="25.83" cy="17.33" r="2.67" />
          <rect x="19.25" y="-1.21" width="1.5" height="22.42" transform="translate(-1.79 15.03) rotate(-39.57)" />
          <circle cx="2.67" cy="17.33" r="2.67" />
          <rect x="-2.71" y="9.25" width="22.42" height="1.5" transform="translate(-4.62 10.18) rotate(-50.44)" />
        </g>
      </g>
    </svg>
  );
};

const InterpolationSmooth: FC<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28.34 20" width={'30px'} height={size} {...rest}>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Icons">
          <circle cx="14.17" cy="2.67" r="2.67" />
          <circle cx="2.67" cy="17.33" r="2.67" />
          <path d="M3.42,17.33H1.92c0-6.46,4.39-15.41,12.64-15.41v1.5C7.29,3.42,3.42,11.5,3.42,17.33Z" />
          <circle cx="25.67" cy="17.33" r="2.67" />
          <path d="M26.42,17.33h-1.5c0-5.83-3.87-13.91-11.14-13.91V1.92C22,1.92,26.42,10.87,26.42,17.33Z" />
        </g>
      </g>
    </svg>
  );
};

const InterpolationStepBefore: FC<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28.34 20" width={'30px'} height={size} {...rest}>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Icons">
          <circle cx="14.17" cy="2.67" r="2.67" />
          <circle cx="2.67" cy="17.33" r="2.67" />
          <circle cx="25.67" cy="17.33" r="2.67" />
          <polygon points="3.42 17.33 1.92 17.33 1.92 1.92 13.78 1.92 13.78 3.42 3.42 3.42 3.42 17.33" />
          <polygon points="25.67 18.08 13.42 18.08 13.42 2.67 14.92 2.67 14.92 16.58 25.67 16.58 25.67 18.08" />
        </g>
      </g>
    </svg>
  );
};

const InterpolationStepAfter: FC<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28.34 20" width={'30px'} height={size} {...rest}>
      <g id="Layer_2" data-name="Layer 2">
        <g id="Icons">
          <circle cx="14.17" cy="2.67" r="2.67" />
          <circle cx="25.67" cy="17.33" r="2.67" />
          <circle cx="2.67" cy="17.33" r="2.67" />
          <polygon points="26.42 17.33 24.92 17.33 24.92 3.42 14.56 3.42 14.56 1.92 26.42 1.92 26.42 17.33" />
          <polygon points="14.92 18.08 2.67 18.08 2.67 16.58 13.42 16.58 13.42 2.67 14.92 2.67 14.92 18.08" />
        </g>
      </g>
    </svg>
  );
};

export const Logs: FC<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg" width={30} height={size} {...rest}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M2.36906 14.2144H8.68657C8.89601 14.2144 9.09687 14.2976 9.24496 14.4457C9.39306 14.5938 9.47626 14.7946 9.47626 15.0041C9.47626 15.2135 9.39306 15.4144 9.24496 15.5625C9.09687 15.7106 8.89601 15.7938 8.68657 15.7938H2.36906C1.74075 15.7938 1.13817 15.5442 0.693883 15.0999C0.249597 14.6556 0 14.053 0 13.4247V2.36906C0 1.74075 0.249597 1.13817 0.693883 0.693883C1.13817 0.249597 1.74075 0 2.36906 0H7.15457L7.37569 0.0789687H7.44676C7.52795 0.116938 7.60259 0.167585 7.66787 0.229009L12.406 4.96714C12.5156 5.07819 12.5898 5.2192 12.6193 5.37239C12.6488 5.52559 12.6323 5.68409 12.5718 5.8279C12.5126 5.97211 12.412 6.09556 12.2827 6.18269C12.1534 6.26982 12.0012 6.31673 11.8453 6.3175H8.68657C8.05825 6.3175 7.45567 6.06791 7.01139 5.62362C6.5671 5.17934 6.3175 4.57675 6.3175 3.94844V1.57938H2.36906C2.15963 1.57938 1.95877 1.66258 1.81067 1.81067C1.66258 1.95877 1.57938 2.15963 1.57938 2.36906V13.4247C1.57938 13.6341 1.66258 13.835 1.81067 13.9831C1.95877 14.1312 2.15963 14.2144 2.36906 14.2144ZM9.94217 4.73813L7.89688 2.69284V3.94844C7.89688 4.15788 7.98008 4.35874 8.12817 4.50683C8.27627 4.65493 8.47713 4.73813 8.68657 4.73813H9.94217ZM3.34761 12.2739C3.31241 12.1914 3.29137 12.1041 3.2854 12.0158L3.2854 8.36528C3.28184 8.18607 3.34961 8.01561 3.47382 7.89141C3.59802 7.76721 3.76848 7.69943 3.94769 7.70299C4.1269 7.70655 4.30019 7.78116 4.42943 7.9104C4.55867 8.03964 4.63327 8.21293 4.63683 8.39214L4.63684 11.2969L5.92266 11.2969C6.01149 11.2983 6.12471 11.3282 6.20753 11.3636C6.27155 11.391 6.31643 11.4391 6.3631 11.4892L6.36311 11.4892C6.37681 11.5039 6.39066 11.5187 6.4052 11.5332C6.46926 11.5973 6.52051 11.6729 6.55597 11.7557C6.59143 11.8386 6.61041 11.9269 6.61181 12.0158C6.61394 12.1046 6.59846 12.1923 6.56626 12.2738C6.53407 12.3553 6.48579 12.4289 6.42422 12.4905C6.36265 12.552 6.289 12.6003 6.20753 12.6325C6.12605 12.6647 6.01151 12.707 5.92266 12.7049H3.97454C3.8858 12.7025 3.79782 12.6813 3.71644 12.6427C3.55164 12.5712 3.41911 12.4387 3.34761 12.2739ZM8.70634 12.793C9.98203 12.793 11.0162 11.6439 11.0162 10.2265C11.0162 8.80904 9.98203 7.65998 8.70634 7.65998C7.43065 7.65998 6.3965 8.80904 6.3965 10.2265C6.3965 11.6439 7.43065 12.793 8.70634 12.793ZM8.73597 11.5453C9.37381 11.5453 9.89089 10.9407 9.89089 10.1949C9.89089 9.44911 9.37381 8.84453 8.73597 8.84453C8.09813 8.84453 7.58105 9.44911 7.58105 10.1949C7.58105 10.9407 8.09813 11.5453 8.73597 11.5453ZM11.3715 10.2265C11.3715 8.80904 12.4056 7.65998 13.6813 7.65998C14.1833 7.65998 14.7779 7.9067 15.0267 8.14005C15.1561 8.2614 15.3112 8.44705 15.3397 8.63063C15.366 8.80006 15.2757 8.96772 15.162 9.08144C15.0708 9.17262 14.938 9.2579 14.8066 9.25912C14.5968 9.26106 14.3907 9.13001 14.3646 9.08144C14.1787 8.93199 13.9535 8.84451 13.7109 8.84451C13.0731 8.84451 12.556 9.44909 12.556 10.1949C12.556 10.9407 13.0731 11.5452 13.7109 11.5452C14.0485 11.5452 14.3522 11.3759 14.5634 11.106L14.6882 10.8582H14.2144C14.1414 10.8571 14.0483 10.8326 13.9802 10.8034C13.9276 10.7809 13.8907 10.7413 13.8523 10.7002C13.8411 10.6882 13.8297 10.676 13.8177 10.664C13.7651 10.6113 13.723 10.5492 13.6938 10.4811C13.6647 10.413 13.6491 10.3404 13.6479 10.2674C13.6462 10.1943 13.6589 10.1222 13.6853 10.0553C13.7118 9.9883 13.7515 9.92777 13.8021 9.87716C13.8527 9.82654 13.9132 9.78686 13.9802 9.7604C14.0472 9.73393 14.1414 9.69913 14.2144 9.70088H15.4247C15.4977 9.70283 15.57 9.72026 15.6369 9.75202C15.7723 9.8108 15.8813 9.91973 15.94 10.0552C15.94 10.0552 15.9629 10.1198 15.9912 10.266C16.0154 10.3912 15.9852 10.5411 15.9549 10.6909C15.9499 10.7159 15.9448 10.741 15.94 10.7659C15.717 11.9243 14.7905 12.793 13.6813 12.793C12.4056 12.793 11.3715 11.6439 11.3715 10.2265Z"
      />
    </svg>
  );
};

const IconNotFound: FC<SvgProps> = ({ size, ...rest }) => {
  return <svg width={size} height={size} {...rest} />;
};

export const customIcons: Record<string, ComponentType<SvgProps>> = {
  'gf-interpolation-linear': InterpolationLinear,
  'gf-interpolation-smooth': InterpolationSmooth,
  'gf-interpolation-step-before': InterpolationStepBefore,
  'gf-interpolation-step-after': InterpolationStepAfter,
  'gf-logs': Logs,
  notFoundDummy: IconNotFound,
};
