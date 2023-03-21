import React, { useEffect, useState } from 'react';
import { IconButton, MenuItem, Menu } from '@mui/material';
import { AccountCircle } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
// import { useLogoutHandler } from '../pkg/hooks';
import ory from '../apis/ory';
import { AxiosError } from 'axios';

function ProfileMenu({logoutUrl}: {logoutUrl: string | undefined}) {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const navigate = useNavigate();

  // const [logoutUrl, setLogoutUrl] = useState<string | undefined>()

  // useEffect(() => {
  //   ory
  //     .createBrowserLogoutFlow()
  //     .then(({ data }) => {
  //       setLogoutUrl(data.logout_token);
  //     })
  //     .catch((err: AxiosError) => {
  //       switch (err.response?.status) {
  //         case 401:
  //           // do nothing, the user is not logged in
  //           return;
  //       }

  //       // Something else happened!
  //       return Promise.reject(err);
  //     });
  // });

  // const logout = useLogoutHandler();

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleProfile = () => {
    handleClose();
    navigate('settings');
  };

  const handleLogout = () => {
    handleClose();
    window.location.href = logoutUrl || '/';
  };

  return (
    <div>
      <IconButton
        size="large"
        aria-label="account of current user"
        aria-controls="menu-appbar"
        aria-haspopup="true"
        onClick={handleMenu}
        color="inherit"
      >
        <AccountCircle />
      </IconButton>
      <Menu
        id="menu-appbar"
        anchorEl={anchorEl}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right'
        }}
        keepMounted
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right'
        }}
        open={Boolean(anchorEl)}
        onClose={handleClose}
      >
        <MenuItem onClick={handleProfile}>Profile</MenuItem>
        <MenuItem onClick={handleLogout}>Logout</MenuItem>
      </Menu>
    </div>
  );
}

export default ProfileMenu;
