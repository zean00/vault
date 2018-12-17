import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import Pretender from 'pretender';

const PERMISSIONS_RESPONSE = {
  data: {
    exact_paths: {
      foo: {
        capabilities: ['read'],
      },
      bar: {
        capabilities: ['create'],
      },
    },
    glob_paths: {
      baz: {
        capabilities: ['read'],
      },
    },
  },
};

module('Unit | Service | permissions', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = new Pretender();
    this.server.get('/v1/sys/internal/ui/resultant-acl', () => {
      return [200, { 'Content-Type': 'application/json' }, JSON.stringify(PERMISSIONS_RESPONSE)];
    });
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('sets paths properly', async function(assert) {
    let service = this.owner.lookup('service:permissions');
    await service.getPaths.perform();
    assert.deepEqual(service.get('paths'), PERMISSIONS_RESPONSE.data);
  });

  test('returns true if a policy includes access to a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    service.set('paths', PERMISSIONS_RESPONSE.data);
    assert.equal(service.hasPermission('foo'), true);
  });

  test('returns false if a policy does not includes access to a path', function(assert) {
    let service = this.owner.lookup('service:permissions');
    assert.equal(service.hasPermission('biz'), false);
  });
});
