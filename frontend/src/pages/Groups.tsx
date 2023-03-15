import { useQuery } from '@apollo/client';
import {
  CircularProgress,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper
} from '@mui/material';
import CreateGroupDialog from '../components/CreateGroupDialog';
import GroupRow from '../components/GroupRow';
import { GroupInfoFragment, useListGroupsQuery } from '../generated/graphql';

function Groups() {
  const { loading, error, data } = useListGroupsQuery();

  return (
    <Box sx={{ padding: 2, backgroundColor: 'white' }}>
      <CreateGroupDialog />
      {loading ? <CircularProgress /> : null}
      {error ? <div>{error.message}</div> : null}
      {data ? (
        <TableContainer
          component={Paper}
          sx={{ padding: 4, marginTop: 3, paddingTop: 1 }}
        >
          <Table aria-label="collapsible table">
            <TableHead>
              <TableRow>
                <TableCell />
                <TableCell>Name</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>
              {data.listGroups?.map((group: GroupInfoFragment ) => (
                  <GroupRow key={group.name} group={group} />
                ))}
            </TableBody>
          </Table>
        </TableContainer>
      ) : null}
    </Box>
  );
}

export default Groups;
