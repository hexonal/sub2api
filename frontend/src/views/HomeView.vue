<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-br from-gray-50 via-primary-50/30 to-gray-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
  >
    <!-- Background Decorations -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div
        class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-primary-400/20 blur-3xl"
      ></div>
      <div
        class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-primary-500/15 blur-3xl"
      ></div>
      <div
        class="absolute left-1/3 top-1/4 h-72 w-72 rounded-full bg-primary-300/10 blur-3xl"
      ></div>
      <div
        class="absolute bottom-1/4 right-1/4 h-64 w-64 rounded-full bg-primary-400/10 blur-3xl"
      ></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(20,184,166,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(20,184,166,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"
      ></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content -->
          <div class="w-full min-w-0 flex-1 text-center lg:text-left">
            <h1
              class="mb-4 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl"
            >
              {{ siteName }}
            </h1>
            <p class="mb-8 text-lg text-gray-600 dark:text-dark-300 md:text-xl">
              {{ siteSubtitle }}
            </p>

            <!-- CTA Button -->
            <div>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-primary-500/30"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
            </div>

            <!-- Install Commands -->
            <div
              class="mx-auto mt-7 w-full max-w-xl rounded-2xl border border-white/70 bg-white/75 p-4 text-left shadow-sm backdrop-blur-sm dark:border-dark-700/70 dark:bg-dark-800/75 lg:mx-0"
            >
              <div class="mb-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div class="flex items-center gap-2">
                  <Icon name="terminal" size="sm" class="text-primary-500" />
                  <span class="text-sm font-semibold text-gray-900 dark:text-white">
                    {{ t('home.install.title') }}
                  </span>
                </div>
                <div
                  class="inline-flex w-full rounded-lg bg-gray-100 p-1 dark:bg-dark-900 sm:w-auto"
                  role="tablist"
                  :aria-label="t('home.install.title')"
                >
                  <button
                    v-for="method in installMethods"
                    :key="method.id"
                    type="button"
                    role="tab"
                    :aria-selected="activeInstallMethod === method.id"
                    @click="activeInstallMethod = method.id"
                    :class="[
                      'flex-1 whitespace-nowrap rounded-md px-3 py-1.5 text-xs font-medium transition-colors sm:flex-none',
                      activeInstallMethod === method.id
                        ? 'bg-white text-primary-600 shadow-sm dark:bg-dark-700 dark:text-primary-300'
                        : 'text-gray-500 hover:text-gray-700 dark:text-dark-300 dark:hover:text-white'
                    ]"
                  >
                    {{ t(method.labelKey) }}
                  </button>
                </div>
              </div>
              <pre
                class="min-h-[52px] overflow-x-auto whitespace-pre-wrap break-all rounded-xl bg-gray-950 px-4 py-3 text-xs leading-6 text-gray-100 shadow-inner sm:text-sm"
              ><code>{{ activeInstallMethodConfig.command }}</code></pre>
            </div>
          </div>

          <!-- Right: Terminal Animation -->
          <div class="flex w-full min-w-0 flex-1 justify-center lg:justify-end">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">terminal</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">entrox</span>
                    <span class="code-flag">login</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># Opening browser authorization...</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">SIGNED IN</span>
                    <span class="code-response">Entrox CLI ready</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Entrox Desktop Downloads -->
        <section
          class="mb-12 overflow-hidden rounded-3xl border border-gray-200/60 bg-white/75 shadow-sm backdrop-blur-sm dark:border-dark-700/60 dark:bg-dark-900/75"
        >
          <div class="grid lg:grid-cols-[0.95fr_1.4fr]">
            <div
              class="relative overflow-hidden bg-gradient-to-br from-gray-950 via-slate-900 to-primary-950 p-8 text-white sm:p-10"
            >
              <div
                class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.05)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.05)_1px,transparent_1px)] bg-[size:56px_56px] opacity-40"
              ></div>
              <div class="relative">
                <div
                  class="mb-8 inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/10 px-4 py-2 text-sm font-semibold text-primary-100 shadow-sm"
                >
                  <Icon name="download" size="sm" />
                  {{ t('home.desktop.badge') }}
                </div>
                <h2 class="mb-5 text-3xl font-bold tracking-normal sm:text-4xl">
                  {{ t('home.desktop.title') }}
                </h2>
                <p class="max-w-md text-base leading-8 text-gray-300">
                  {{ t('home.desktop.description') }}
                </p>
              </div>
            </div>

            <div class="grid gap-4 p-5 sm:p-8 md:grid-cols-2 xl:grid-cols-4">
              <article
                v-for="platform in desktopPlatforms"
                :key="platform.id"
                class="flex min-h-[240px] flex-col rounded-2xl border border-gray-200/80 bg-white p-6 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/80 dark:bg-dark-800"
              >
                <div
                  :class="[
                    'mb-6 flex h-16 w-16 items-center justify-center rounded-2xl text-xl font-bold text-white shadow-lg',
                    platform.iconClass
                  ]"
                >
                  {{ platform.badge }}
                </div>
                <h3 class="mb-3 text-lg font-semibold text-gray-900 dark:text-white">
                  {{ t(platform.titleKey) }}
                </h3>
                <p class="mb-6 text-sm leading-6 text-gray-500 dark:text-dark-300">
                  {{ t(platform.descriptionKey) }}
                </p>
                <a
                  v-if="platform.downloadUrl"
                  :href="platform.downloadUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="mt-auto inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-gray-900 px-4 text-sm font-semibold text-white transition-colors hover:bg-gray-800 dark:bg-primary-600 dark:hover:bg-primary-500"
                >
                  <Icon name="download" size="sm" :stroke-width="2" />
                  {{ t('home.desktop.download') }}
                </a>
                <button
                  v-else
                  type="button"
                  disabled
                  :title="t('home.desktop.downloadPending')"
                  class="mt-auto inline-flex h-10 cursor-not-allowed items-center justify-center gap-2 rounded-lg border border-gray-200 bg-gray-100 px-4 text-sm font-semibold text-gray-400 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-500"
                >
                  <Icon name="download" size="sm" :stroke-width="2" />
                  {{ t('home.desktop.download') }}
                </button>
              </article>
            </div>
          </div>
        </section>

        <!-- Personal Subscription Pricing -->
        <section class="mb-12">
          <div class="mb-7 text-center">
            <div
              class="mb-4 inline-flex items-center gap-2 rounded-full border border-primary-200/70 bg-white/80 px-4 py-2 text-sm font-semibold text-primary-700 shadow-sm backdrop-blur-sm dark:border-primary-800/60 dark:bg-dark-800/80 dark:text-primary-300"
            >
              <Icon name="creditCard" size="sm" />
              {{ t('home.pricing.badge') }}
            </div>
            <h2 class="mb-3 text-3xl font-bold tracking-normal text-gray-900 dark:text-white">
              {{ t('home.pricing.title') }}
            </h2>
            <p class="mx-auto max-w-2xl text-base leading-7 text-gray-600 dark:text-dark-300">
              {{ t('home.pricing.description') }}
            </p>
          </div>

          <div class="grid gap-5 sm:grid-cols-2 xl:grid-cols-4">
            <article
              v-for="plan in pricingPlans"
              :key="plan.id"
              :class="[
                'relative flex min-h-[360px] flex-col rounded-2xl border bg-white/80 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 dark:bg-dark-800/80',
                plan.featured
                  ? 'border-primary-300 shadow-primary-500/15 ring-1 ring-primary-300/70 dark:border-primary-700 dark:ring-primary-700/80'
                  : 'border-gray-200/70 hover:shadow-xl hover:shadow-gray-900/5 dark:border-dark-700/70'
              ]"
            >
              <div
                v-if="plan.featured"
                class="absolute right-5 top-5 inline-flex items-center gap-1.5 rounded-full bg-primary-500 px-3 py-1 text-xs font-semibold text-white shadow-sm"
              >
                <Icon name="badge" size="xs" :stroke-width="2" />
                {{ t('home.pricing.recommended') }}
              </div>

              <div class="mb-6 pr-24">
                <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ t(plan.nameKey) }}
                </h3>
                <p class="mt-2 text-sm leading-6 text-gray-500 dark:text-dark-300">
                  {{ t(plan.subtitleKey) }}
                </p>
              </div>

              <div class="mb-6 flex items-end gap-2">
                <span
                  v-if="plan.price"
                  class="text-4xl font-bold tracking-normal text-gray-900 dark:text-white"
                >
                  ¥{{ plan.price }}
                </span>
                <span
                  v-else
                  class="text-3xl font-bold tracking-normal text-gray-900 dark:text-white"
                >
                  {{ t(plan.priceLabelKey) }}
                </span>
                <span
                  v-if="plan.price"
                  class="pb-1 text-sm font-medium text-gray-500 dark:text-dark-300"
                >
                  {{ t('home.pricing.month') }}
                </span>
              </div>

              <ul class="mb-7 flex flex-1 flex-col gap-3">
                <li
                  v-for="featureKey in plan.featureKeys"
                  :key="featureKey"
                  class="flex items-start gap-2.5 text-sm leading-6 text-gray-600 dark:text-dark-300"
                >
                  <Icon
                    name="checkCircle"
                    size="sm"
                    class="mt-1 shrink-0 text-primary-500"
                    :stroke-width="2"
                  />
                  <span>{{ t(featureKey) }}</span>
                </li>
              </ul>

              <a
                v-if="plan.ctaHref"
                :href="plan.ctaHref"
                :class="pricingPlanCtaClass(plan.featured)"
              >
                {{ t(plan.ctaKey) }}
              </a>
              <router-link
                v-else
                :to="isAuthenticated ? dashboardPath : '/login'"
                :class="pricingPlanCtaClass(plan.featured)"
              >
                {{ t(plan.ctaKey) }}
              </router-link>
            </article>
          </div>

          <p class="mt-5 text-center text-sm leading-6 text-gray-500 dark:text-dark-400">
            {{ t('home.pricing.note') }}
          </p>
        </section>

        <!-- Features Grid -->
        <div class="mb-12 grid gap-6 md:grid-cols-3">
          <!-- Feature 1: Unified Gateway -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 shadow-lg shadow-blue-500/30 transition-transform group-hover:scale-110"
            >
              <Icon name="server" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.unifiedGateway') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.unifiedGatewayDesc') }}
            </p>
          </div>

          <!-- Feature 2: Account Pool -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-primary-500 to-primary-600 shadow-lg shadow-primary-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.multiAccount') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.multiAccountDesc') }}
            </p>
          </div>

          <!-- Feature 3: Billing & Quota -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-primary-500/10 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 shadow-lg shadow-purple-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.balanceQuota') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.balanceQuotaDesc') }}
            </p>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="mb-16 flex flex-wrap items-center justify-center gap-4">
          <!-- Claude - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-orange-400 to-orange-500"
            >
              <span class="text-xs font-bold text-white">C</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.claude') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- GPT - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-green-500 to-green-600"
            >
              <span class="text-xs font-bold text-white">G</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">GPT</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Gemini - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-blue-600"
            >
              <span class="text-xs font-bold text-white">G</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.gemini') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- More - Coming Soon -->
          <div
            class="flex items-center gap-2 rounded-xl border border-gray-200/50 bg-white/40 px-5 py-3 opacity-60 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/40"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-gray-500 to-gray-600"
            >
              <span class="text-xs font-bold text-white">+</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.more') }}</span>
            <span
              class="rounded bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700 dark:text-dark-400"
              >{{ t('home.providers.soon') }}</span
            >
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { getEntroxInstallMethods, isCnInstallLocale, type EntroxInstallMethodId } from '@/utils/entroxInstall'

const { t, locale } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

type DesktopPlatformId = 'mac-arm' | 'mac-intel' | 'windows' | 'linux'
type PricingPlanId = 'pro' | 'plus' | 'ultra' | 'enterprise'

interface DesktopPlatform {
  id: DesktopPlatformId
  badge: string
  titleKey: string
  descriptionKey: string
  iconClass: string
  downloadUrl: string
}

interface PricingPlan {
  id: PricingPlanId
  price: string
  priceLabelKey: string
  nameKey: string
  subtitleKey: string
  featureKeys: string[]
  ctaKey: string
  ctaHref?: string
  featured: boolean
}

const defaultSiteName = 'Entrox Studio'
const defaultSiteSubtitle =
  'Entrox 把模型网关、CLI 与 Desktop 工作台连成一套 AI 编程系统：一键登录同步 Claude/GPT/Gemini 等模型资源，在桌面同时管理多个项目与 Agent 会话，实时查看流式输出、排队提示词、切换模型/后端并接入 MCP 工具；在终端用 entrox 快速启动同一套能力，让日常改动、自动化任务和长程编码都保持连续、清晰、可控。'

const normalizeEntroxBranding = (value: string | undefined, fallback: string) => {
  const trimmed = value?.trim()
  if (
    !trimmed ||
    trimmed === 'Sub2API' ||
    trimmed === 'AI API Gateway Platform' ||
    trimmed === 'Subscription to API Conversion Platform'
  ) {
    return fallback
  }
  return trimmed
    .replace(/\bEntro Studio\b/g, 'Entrox Studio')
    .replace(/\bEntro Desktop\b/g, 'Entrox Desktop')
    .replace(/\bCli\b/g, 'CLI')
}

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() =>
  normalizeEntroxBranding(appStore.cachedPublicSettings?.site_name || appStore.siteName, defaultSiteName)
)
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const siteSubtitle = computed(() =>
  normalizeEntroxBranding(appStore.cachedPublicSettings?.site_subtitle, defaultSiteSubtitle)
)
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

const activeInstallMethod = ref<EntroxInstallMethodId>('script')
const installMethods = computed(() => getEntroxInstallMethods(isCnInstallLocale(locale.value)))
const activeInstallMethodConfig = computed(
  () => installMethods.value.find((method) => method.id === activeInstallMethod.value) || installMethods.value[0]
)
const desktopDownloadBaseUrl = 'https://entrox-download.996icu.wiki/entrox-desktop'
const desktopPlatforms: DesktopPlatform[] = [
  {
    id: 'mac-arm',
    badge: 'M',
    titleKey: 'home.desktop.platforms.macArm.title',
    descriptionKey: 'home.desktop.platforms.macArm.description',
    iconClass: 'bg-gradient-to-br from-gray-900 to-gray-700 shadow-gray-900/20',
    downloadUrl: `${desktopDownloadBaseUrl}/entrox-desktop-mac-arm64.dmg`
  },
  {
    id: 'mac-intel',
    badge: 'Intel',
    titleKey: 'home.desktop.platforms.macIntel.title',
    descriptionKey: 'home.desktop.platforms.macIntel.description',
    iconClass: 'bg-gradient-to-br from-slate-700 to-slate-500 shadow-slate-700/20',
    downloadUrl: `${desktopDownloadBaseUrl}/entrox-desktop-mac-x64.dmg`
  },
  {
    id: 'windows',
    badge: 'Win',
    titleKey: 'home.desktop.platforms.windows.title',
    descriptionKey: 'home.desktop.platforms.windows.description',
    iconClass: 'bg-gradient-to-br from-blue-500 to-cyan-500 shadow-blue-500/25',
    downloadUrl: `${desktopDownloadBaseUrl}/entrox-desktop-0.5.14.exe`
  },
  {
    id: 'linux',
    badge: 'Linux',
    titleKey: 'home.desktop.platforms.linux.title',
    descriptionKey: 'home.desktop.platforms.linux.description',
    iconClass: 'bg-gradient-to-br from-orange-500 to-amber-500 shadow-orange-500/25',
    downloadUrl: `${desktopDownloadBaseUrl}/entrox-desktop-0.5.14.deb`
  }
]
const pricingPlans: PricingPlan[] = [
  {
    id: 'pro',
    price: '59',
    priceLabelKey: '',
    nameKey: 'home.pricing.plans.pro.name',
    subtitleKey: 'home.pricing.plans.pro.subtitle',
    featureKeys: [
      'home.pricing.plans.pro.usage',
      'home.pricing.plans.pro.concurrency',
      'home.pricing.plans.pro.queue'
    ],
    ctaKey: 'home.pricing.cta',
    featured: false
  },
  {
    id: 'plus',
    price: '159',
    priceLabelKey: '',
    nameKey: 'home.pricing.plans.plus.name',
    subtitleKey: 'home.pricing.plans.plus.subtitle',
    featureKeys: [
      'home.pricing.plans.plus.usage',
      'home.pricing.plans.plus.concurrency',
      'home.pricing.plans.plus.queue'
    ],
    ctaKey: 'home.pricing.cta',
    featured: true
  },
  {
    id: 'ultra',
    price: '399',
    priceLabelKey: '',
    nameKey: 'home.pricing.plans.ultra.name',
    subtitleKey: 'home.pricing.plans.ultra.subtitle',
    featureKeys: [
      'home.pricing.plans.ultra.usage',
      'home.pricing.plans.ultra.concurrency',
      'home.pricing.plans.ultra.queue'
    ],
    ctaKey: 'home.pricing.cta',
    featured: false
  },
  {
    id: 'enterprise',
    price: '',
    priceLabelKey: 'home.pricing.plans.enterprise.price',
    nameKey: 'home.pricing.plans.enterprise.name',
    subtitleKey: 'home.pricing.plans.enterprise.subtitle',
    featureKeys: [
      'home.pricing.plans.enterprise.usage',
      'home.pricing.plans.enterprise.concurrency',
      'home.pricing.plans.enterprise.support'
    ],
    ctaKey: 'home.pricing.plans.enterprise.cta',
    ctaHref: 'mailto:ailun@matchtrio.com',
    featured: false
  }
]

const pricingPlanCtaClass = (featured: boolean) => [
  'mt-auto inline-flex h-11 items-center justify-center rounded-lg px-4 text-sm font-semibold transition-colors',
  featured
    ? 'bg-primary-600 text-white shadow-lg shadow-primary-500/20 hover:bg-primary-500'
    : 'border border-gray-200 bg-white text-gray-900 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-900 dark:text-white dark:hover:bg-dark-700'
]

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* Terminal Container */
.terminal-container {
  position: relative;
  display: inline-block;
  width: min(420px, 100%);
  max-width: 100%;
}

/* Terminal Window */
.terminal-window {
  width: 100%;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.5s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #14b8a6;
}
.code-comment {
  color: #64748b;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fbbf24;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:deep(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(20, 184, 166, 0.2),
    0 0 40px rgba(20, 184, 166, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}
</style>
