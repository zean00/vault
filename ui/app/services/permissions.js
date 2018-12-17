import Service, { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

// decide what happens when no policies are found
// because requesting the root token's ACL policy always returns nil
// https://github.com/hashicorp/vault/pull/4386/files#diff-a145099c8e3d917858619dcfb2ae09b9R3464

export default Service.extend({
  paths: null,
  store: service(),

  getPaths: task(function*() {
    if (this.get('paths')) {
      return;
    }
    let resp = yield this.get('store')
      .adapterFor('permissions')
      .query();
    this.setPaths(resp);
    return;
  }),

  setPaths(resp) {
    this.set('paths', resp.data);
  },

  hasPermission(pathName) {
    if (this.get('paths')) {
      const paths = this.get('paths');
      return Object.keys(paths).some(pathType => {
        return paths[pathType].hasOwnProperty(pathName);
      });
    }
    return false;
  },
});
