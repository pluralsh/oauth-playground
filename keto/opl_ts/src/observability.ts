import { Namespace, SubjectSet, Context } from '@ory/keto-namespace-types'

class Admin implements Namespace {}

class Organization implements Namespace {
  related: {
    admins: Admin[]
  }
}

class User implements Namespace {
  related: {
    admins: SubjectSet<Organization, 'admins'>[]
  }
  permits = {
    create: (ctx: Context) => this.related.admins.includes(ctx.subject),
    edit: (ctx: Context) => this.related.admins.includes(ctx.subject),
    delete: (ctx: Context) => this.related.admins.includes(ctx.subject),
  }
}

class Group implements Namespace {
  related: {
    admins: SubjectSet<Organization, 'admins'>[]
    members: (User | Group)[]
  }
  permits = {
    create: (ctx: Context) => this.related.admins.includes(ctx.subject),
    edit: (ctx: Context) => this.related.admins.includes(ctx.subject),
    delete: (ctx: Context) => this.related.admins.includes(ctx.subject),
  }
}

class OAuth2Client implements Namespace {
  related: {
    admins: SubjectSet<Organization, 'admins'>[]
    loginBindings: (User | SubjectSet<Group, 'members'>)[]
  }
  permits = {
    login: (ctx: Context) =>
      this.related.loginBindings.includes(ctx.subject) ||
      this.related.admins.includes(ctx.subject),
    edit: (ctx: Context) => this.related.admins.includes(ctx.subject),
    delete: (ctx: Context) => this.related.admins.includes(ctx.subject),
    create: (ctx: Context) => this.related.admins.includes(ctx.subject),
  }
}

class ObservabilityTenant implements Namespace {
  related: {
    admins: SubjectSet<Organization, 'admins'>[]
    viewers: (User | SubjectSet<Group, 'members'> | OAuth2Client)[]
    editors: (User | SubjectSet<Group, 'members'>)[]
  }
  permits = {
    view: (ctx: Context) =>
      this.related.viewers.includes(ctx.subject) ||
      this.related.editors.includes(ctx.subject) ||
      this.related.admins.includes(ctx.subject),
    edit: (ctx: Context) =>
      this.related.editors.includes(ctx.subject) ||
      this.related.admins.includes(ctx.subject),
    delete: (ctx: Context) => this.related.admins.includes(ctx.subject),
    create: (ctx: Context) => this.related.admins.includes(ctx.subject),
  }
}
