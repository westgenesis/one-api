import React, { useEffect, useState } from 'react';
import { API, showError, showSuccess } from '../helpers';
import { Form, Button, Label, Table, Popup } from 'semantic-ui-react';

const GroupsTable = () => {
  const [loading, setLoading] = useState(true);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [addLoading, setAddLoading] = useState(false);
  const [groups, setGroups] = useState([]);
  const [addGroup, setAddGroup] = useState('');

  const loadGroups = async () => {
    const res = await API.get('/api/group');
    const { success, message, data } = res.data;
    if (success) {
      setGroups(data);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const handleDelete = async (name) => {
    setDeleteLoading(true);
    const res = await API.post(`/api/group/delete`, {
      name,
    });
    const { success, message } = res.data;
    if (success) {
      showSuccess('删除成功!');
      loadGroups();
    } else {
      showError(message);
    }
    setDeleteLoading(false);
  };

  const handleAddGroup = async () => {
    if (!addGroup) {
      showError('请输入分组名称');
      return;
    }
    setAddLoading(true);
    const res = await API.post('/api/group/create', {
      name: addGroup,
      ratio: 1,
    });
    const { success, message } = res.data;
    if (success) {
      showSuccess('添加成功!');
      setAddGroup('');
      loadGroups();
    } else {
      showError(message);
    }
    setAddLoading(false);
  };

  useEffect(() => {
    loadGroups()
      .then()
      .catch((reason) => {
        showError(reason);
      });
  }, []);

  return (
    <>
      <Form>
        <div
          style={{
            display: 'flex',
            gap: '10px',
          }}
        >
          <Form.Input
            style={{
              flex: 1,
              minWidth: '400px',
            }}
            inline
            label='分组名称'
            placeholder='请输入分组名称'
            value={addGroup}
            onChange={(e) => setAddGroup(e.target.value)}
          />
          <Form.Button loading={addLoading} onClick={handleAddGroup}>
            添加
          </Form.Button>
        </div>
      </Form>
      <Table celled>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>分组名称</Table.HeaderCell>
            <Table.HeaderCell width={2} textAlign='center'>
              操作
            </Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {groups.map((group) => (
            <Table.Row>
              <Table.Cell>
                <Label>{group}</Label>
              </Table.Cell>
              <Table.Cell textAlign='center'>
                <Popup
                  trigger={
                    <Button size='tiny' negative loading={deleteLoading}>
                      删除
                    </Button>
                  }
                  on='click'
                  flowing
                  hoverable
                >
                  <Button
                    negative
                    size={'tiny'}
                    onClick={() => handleDelete(group)}
                  >
                    删除分组 {group}
                  </Button>
                </Popup>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table>
    </>
  );
};

export default GroupsTable;
