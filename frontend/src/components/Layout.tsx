import { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';
import { createStyles, makeStyles } from '@mui/styles';
import { Theme, styled, useTheme, CSSObject } from '@mui/material/styles';
import {
  Box,
  Toolbar,
  List,
  Divider,
  IconButton,
  ListSubheader,
  Typography,
  CssBaseline,
  ListItemIcon,
  ListItemText,
  ListItemButton
} from '@mui/material';
import MuiDrawer from '@mui/material/Drawer';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import {
  Menu,
  ChevronLeft,
  ChevronRight,
  Home,
  PieChart,
  Folder,
  Book,
  Storage,
  Assessment,
  Person,
  Group
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import ProfileMenu from './ProfileMenu';
import { sdk, sdkError } from "../apis/ory"
import { Session } from '@ory/client';
// import SelectNamespaceMenu from './SelectNamespaceMenu';
// import { ME } from '../graphql/queries';
import { useQuery } from '@apollo/client';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    pages: {
      background: '#f9f9f9',
      width: '100%'
    }
  })
);

const drawerWidth = 240;

const openedMixin = (theme: Theme): CSSObject => ({
  width: drawerWidth,
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.enteringScreen
  }),
  overflowX: 'hidden'
});

const closedMixin = (theme: Theme): CSSObject => ({
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen
  }),
  overflowX: 'hidden',
  width: `calc(${theme.spacing(7)} + 1px)`,
  [theme.breakpoints.up('sm')]: {
    width: `calc(${theme.spacing(9)} + 1px)`
  }
});

const DrawerHeader = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: prop => prop !== 'open'
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  })
}));

const Drawer = styled(MuiDrawer, {
  shouldForwardProp: prop => prop !== 'open'
})(({ theme, open }) => ({
  width: drawerWidth,
  flexShrink: 0,
  whiteSpace: 'nowrap',
  boxSizing: 'border-box',
  ...(open && {
    ...openedMixin(theme),
    '& .MuiDrawer-paper': openedMixin(theme)
  }),
  ...(!open && {
    ...closedMixin(theme),
    '& .MuiDrawer-paper': closedMixin(theme)
  })
}));

export default function Layout() {
  const classes = useStyles();
  const navigate = useNavigate();
  const location = useLocation();

  // const { data, loading, error } = useQuery(ME);

  const theme = useTheme();
  const [open, setOpen] = useState(false);

  const [session, setSession] = useState<Session | undefined>();
  const [logoutUrl, setLogoutUrl] = useState<string>()

  const createLogoutFlow = () => {
    // here we create a new logout URL which we can use to log the user out
    sdk
      .createBrowserLogoutFlow(undefined, {
        params: {
          return_url: "/",
        },
      })
      .then(({ data }) => setLogoutUrl(data.logout_url))
      .catch(sdkErrorHandler)
  }

  const sdkErrorHandler = sdkError(undefined, undefined, "/login")

  useEffect(() => {
    sdk
      .toSession()
      .then(({ data: session }) => {
        setSession(session);
        createLogoutFlow()
      })
      .catch(sdkErrorHandler)
      .catch((error) => {
        // Handle all other errors like error.message "network error" if Kratos can not be connected etc.
        if (error.message) {
          return navigate(`/error?error=${encodeURIComponent(error.message)}`, {
            replace: true,
          })
        }

        // Just stringify error and print all data
        navigate(`/error?error=${encodeURIComponent(JSON.stringify(error))}`, {
          replace: true,
        })
      });
  }, []);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  const mainMenuItems = [
    {
      text: 'Home',
      icon: <Home />,
      path: '/'
    },
    {
      text: 'Dashboards',
      icon: <PieChart />,
      path: '/create'
    },
    {
      text: 'Projects',
      icon: <Folder />,
      path: '/empty'
    }
  ];

  const workspaceMenuItems = [
    {
      text: 'Workspaces',
      icon: <Book />,
      path: '/workspaces'
    },
    {
      text: 'Storage',
      icon: <Storage />,
      path: '/storage'
    },
    {
      text: 'Tensorboards',
      icon: <Assessment />,
      path: '/empty'
    }
  ];

  const usersMenuItems = [
    {
      text: 'Users',
      icon: <Person />,
      path: '/users'
    },
    {
      text: 'Groups',
      icon: <Group />,
      path: '/groups'
    }
  ];

// TODO: add back admin check

  return (
    <div>
      <Box sx={{ display: 'flex' }}>
        <CssBaseline />
        <AppBar position="fixed" open={open}>
          <Toolbar sx={{ display: 'flex', justifyContent: 'space-between' }}>
            <Box sx={{ display: 'flex', flexDirection: 'row' }}>
              <IconButton
                color="inherit"
                aria-label="open drawer"
                onClick={handleDrawerOpen}
                edge="start"
                sx={{
                  marginRight: 2,
                  ...(open && { display: 'none' })
                }}
              >
                <Menu />
              </IconButton>
              <Typography variant="h6" noWrap component="div" mt={0.5}>
                Kubricks
              </Typography>
              <Box marginLeft={2}>
                {/* <SelectNamespaceMenu /> */}
              </Box>
            </Box>
            <ProfileMenu logoutUrl={logoutUrl}/>
          </Toolbar>
        </AppBar>
        <Drawer variant="permanent" open={open}>
          <DrawerHeader>
            <IconButton onClick={handleDrawerClose}>
              {theme.direction === 'rtl' ? <ChevronRight /> : <ChevronLeft />}
            </IconButton>
          </DrawerHeader>
          <Divider />
          <List>
            <ListSubheader>Main</ListSubheader>
            {mainMenuItems.map(item => (
              <ListItemButton
                key={item.text}
                onClick={() => navigate(item.path)}
                selected={location.pathname === item.path ? true : false}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.text}></ListItemText>
              </ListItemButton>
            ))}
          </List>
          <Divider />
          <List>
            <ListSubheader>Env</ListSubheader>
            {workspaceMenuItems.map(item => (
              <ListItemButton
                key={item.text}
                onClick={() => navigate(item.path)}
                selected={location.pathname === item.path ? true : false}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.text}></ListItemText>
              </ListItemButton>
            ))}
          </List>
          {/* {data && data.me && data.me.admin ? ( */}
            <Box>
              <Divider />
              <List>
                <ListSubheader>Env</ListSubheader>
                {usersMenuItems.map(item => (
                  <ListItemButton
                    key={item.text}
                    onClick={() => navigate(item.path)}
                    selected={location.pathname === item.path ? true : false}
                  >
                    <ListItemIcon>{item.icon}</ListItemIcon>
                    <ListItemText primary={item.text}></ListItemText>
                  </ListItemButton>
                ))}
              </List>
            </Box>
          {/* ) : null} */}
        </Drawer>
        <Box component="main" sx={{ flexGrow: 1 }}>
          <DrawerHeader />
          <div>
            <Outlet />
          </div>
        </Box>
      </Box>
    </div>
  );
}
