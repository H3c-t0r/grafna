import React, { FunctionComponent } from 'react';
import { SvgProps } from '../types';

export const PerconaClusterNetwork: FunctionComponent<SvgProps> = ({ size, ...rest }) => {
  return (
    <svg width={size} height={size} {...rest} viewBox="0 0 23 23" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M19.0374 13.0366V9.96332C19.3008 9.90668 19.5575 9.81334 19.7941 9.67666H19.7974C21.015 8.97178 21.4468 7.40818 20.7375 6.16665C19.8439 4.62697 17.7047 4.46883 16.5975 5.74331L13.9375 4.20998C14.4943 2.53956 13.2292 0.833313 11.5008 0.833313C9.74562 0.833313 8.51027 2.5584 9.0608 4.20998L6.40414 5.74664C5.22781 4.41796 3.12698 4.67417 2.26083 6.16665C1.43127 7.61735 2.16266 9.56505 3.96745 9.96665V13.0366C0.997887 13.6307 1.38358 18.12 4.48747 18.12C5.19572 18.12 5.90226 17.821 6.4008 17.2566L9.0608 18.7933C8.50807 20.4302 9.73952 22.1666 11.5008 22.1666C13.1589 22.1666 14.4973 20.5746 13.9375 18.7933L16.5975 17.2566C17.7585 18.5683 19.8627 18.3371 20.7375 16.84C20.7375 16.8366 20.7375 16.8366 20.7408 16.8333C21.5815 15.3502 20.8005 13.4161 19.0374 13.0366ZM16.2208 16.7033V16.7066L13.6475 18.19C12.6513 16.617 10.3413 16.7338 9.3508 18.19L6.78081 16.7066C7.00773 16.2523 7.10701 15.6689 7.01417 15.11C6.94874 15.0445 6.70491 13.0696 4.63412 12.9867V10.0166C6.23211 10.0166 7.66579 8.06995 6.78081 6.29999L9.35413 4.81334C9.36529 4.83566 9.52852 5.05628 9.55412 5.07335C10.2403 5.87945 11.3026 6.1337 12.2108 5.86664C12.2208 5.86664 12.2308 5.85997 12.2408 5.85997C12.756 5.68831 13.1361 5.41265 13.1675 5.34996C13.3508 5.1933 13.5142 5.01662 13.6475 4.81334L16.2208 6.29665C15.5381 7.74759 16.131 9.69945 18.3708 10.0166V12.9867C16.5823 13.0716 15.4079 14.917 16.2208 16.7033Z"
        fill="#9FA7B3"
      />
      <path
        d="M11.6675 7.3366C11.3676 7.16244 11.411 7.36201 7.97752 9.27658C7.67982 9.44946 7.87184 9.487 7.81085 13.4399C7.81085 13.5566 7.87416 13.6666 7.97752 13.7266C11.23 15.5287 11.3544 15.8544 11.6675 15.6666L15.0242 13.7266C15.3253 13.5517 15.1302 13.4894 15.1909 9.56328C15.1909 9.21766 15.0143 9.34088 11.6675 7.3366Z"
        fill="#9FA7B3"
      />
    </svg>
  );
};
