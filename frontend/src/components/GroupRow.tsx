import { useState } from 'react';
import { useMutation } from '@apollo/client';
import {
  Collapse,
  IconButton,
  Box,
  Typography,
  Table,
  TableBody,
  TableHead,
  TableRow,
  TableCell
} from '@mui/material';
import {
  KeyboardArrowDown,
  KeyboardArrowUp,
  Delete
} from '@mui/icons-material';
import { namedOperations, GroupInfoFragment, useDeleteGroupMutation, User } from '../generated/graphql';

function GroupRow({ group }: { group: GroupInfoFragment }) {
  const [open, setOpen] = useState(false);
  const [deleteGroup, { loading, error }] = useDeleteGroupMutation({
    variables: {
      name: group.name || ''
    },
    refetchQueries: [namedOperations.Query.ListGroups]
  });

  return (
    <>
      <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {group.name}
        </TableCell>
        <TableCell align="right">
          <Box display="flex" flexDirection="row" justifyContent="end">
            <IconButton size="small" onClick={() => deleteGroup()}>
              <Delete />
            </IconButton>
          </Box>
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Users
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell>Email</TableCell>
                    <TableCell>Name</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {group.members
                    ? group.members.map((user: Partial<User>) => (
                        <TableRow key={user.email}>
                          <TableCell component="th" scope="row">
                            {user.email}
                          </TableCell>
                          <TableCell component="th" scope="row">
                            {user.name?.first} {user.name?.last}
                          </TableCell>
                        </TableRow>
                      ))
                    : null}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  );
}

export default GroupRow;
