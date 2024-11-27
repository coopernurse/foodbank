const m = window.m;
import i18n from '../services/i18n.js';
import HouseholdService from '../services/HouseholdService.js';

const Signup = {
  oninit: vnode => {
    vnode.state.household = {
      head: {},
      members: []
    };
    vnode.state.errors = {};
    vnode.state.lang = i18n.currentLang;
  },

  view: vnode => {
    const { household, errors, lang } = vnode.state;

    const renderField = (labelKey, name, type = 'text', value = '') => {
      return m('div', { class: 'mb-4' }, [
        m('label', { class: 'block text-gray-700 text-sm font-bold mb-2', for: name }, i18n.t(labelKey)),
        m('input', {
          class: `shadow appearance-none border ${errors[name] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
          type,
          name,
          value,
          oninput: function (ev) { household.head[name] = ev.target.value; }
        }),
        errors[name] && m('p', { class: 'text-red-500 text-xs italic' }, errors[name])
      ]);
    };

    const renderSelect = (labelKey, name, options) => {
      return m('div', { class: 'mb-4' }, [
        m('label', { class: 'block text-gray-700 text-sm font-bold mb-2', for: name }, i18n.t(labelKey)),
        m('select', {
          class: `shadow appearance-none border ${errors[name] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
          name,
          oninput: function (ev) { household.head[name] = ev.target.value; }
        }, options.map(option => m('option', { value: option.value }, option.label)))
      ]);
    };

    const renderMemberFields = (memberIndex) => {
      return m('div', { class: 'mb-4' }, [
        renderField(`misc.firstname`, `member${memberIndex}FirstName`),
        renderField(`misc.lastname`, `member${memberIndex}LastName`),
        renderSelect(`misc.gender`, `member${memberIndex}Gender`, [
          { value: 'male', label: i18n.t('misc.male') },
          { value: 'female', label: i18n.t('misc.female') },
          { value: 'optout', label: i18n.t('misc.prefernottosay') }
        ]),
        renderSelect(`misc.dob`, `member${memberIndex}Dob`, [
          { value: '01', label: '01' },
          { value: '02', label: '02' },
          // Add more options for days, months, and years
        ])
      ]);
    };

    return m('div', { class: 'container mx-auto p-4' }, [
      m('h1', { class: 'text-2xl font-bold mb-4' }, i18n.t('signup.title')),
      m('p', { class: 'mb-4' }, i18n.t('signup.intro')),
      m('div', { class: 'mb-4' }, [
        m('a', { href: '#', onclick: () => i18n.setLanguage('en'), class: 'mr-4' }, 'English'),
        m('a', { href: '#', onclick: () => i18n.setLanguage('es') }, 'Español')
      ]),
      m('form', {
        onsubmit: async e => {
          e.preventDefault();
          try {
            await HouseholdService.createHousehold(household);
            m.route.set('/success');
          } catch (err) {
            vnode.state.errors = { 'form': i18n.t('misc.error') };
          }
        },
        class: 'bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4'
      }, [
        m('h2', { class: 'text-xl font-bold mb-4' }, i18n.t('signup.hoh')),
        renderField(`misc.firstname`, 'firstName'),
        renderField(`misc.lastname`, 'lastName'),
        renderField(`misc.address`, 'address'),
        renderField(`misc.city`, 'city'),
        renderField(`misc.zipcode`, 'zipcode'),
        renderField(`misc.email`, 'email'),
        renderField(`misc.phone`, 'phone'),
        renderSelect(`misc.gender`, 'gender', [
          { value: 'male', label: i18n.t('misc.male') },
          { value: 'female', label: i18n.t('misc.female') },
          { value: 'optout', label: i18n.t('misc.prefernottosay') }
        ]),
        renderSelect(`misc.dob`, 'dob', [
          { value: '01', label: '01' },
          { value: '02', label: '02' },
          // Add more options for days, months, and years
        ]),
        m('h2', { class: 'text-xl font-bold mt-8 mb-4' }, i18n.t('signup.othermembers')),
        Array.from({ length: 5 }).map((_, i) => renderMemberFields(i)),
        m('div', { class: 'flex items-center justify-between' }, [
          m('button', {
            type: 'submit',
            class: 'bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline'
          }, i18n.t('misc.submit'))
        ])
      ])
    ]);
  }
};

export default Signup;
