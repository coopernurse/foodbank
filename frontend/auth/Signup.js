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
        m('div', { class: 'mb-4' }, [
          m('label', { class: 'block text-gray-700 text-sm font-bold mb-2', for: `member${memberIndex}Dob` }, i18n.t('misc.dob')),
          m('div', { class: 'flex' }, [
            m('select', {
              class: `shadow appearance-none border ${errors[`member${memberIndex}DobMonth`] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: `member${memberIndex}DobMonth`,
              oninput: function (ev) { household.members[memberIndex]['dobMonth'] = ev.target.value; }
            }, [
              { value: '01', label: '01' },
              { value: '02', label: '02' },
              // Add more options for months
            ].map(option => m('option', { value: option.value }, option.label))),
            m('select', {
              class: `shadow appearance-none border ${errors[`member${memberIndex}DobDay`] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: `member${memberIndex}DobDay`,
              oninput: function (ev) { household.members[memberIndex]['dobDay'] = ev.target.value; }
            }, [
              { value: '01', label: '01' },
              { value: '02', label: '02' },
              // Add more options for days
            ].map(option => m('option', { value: option.value }, option.label))),
            m('select', {
              class: `shadow appearance-none border ${errors[`member${memberIndex}DobYear`] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: `member${memberIndex}DobYear`,
              oninput: function (ev) { household.members[memberIndex]['dobYear'] = ev.target.value; }
            }, [
              { value: '2000', label: '2000' },
              { value: '2001', label: '2001' },
              // Add more options for years
            ].map(option => m('option', { value: option.value }, option.label)))
          ])
        ]),
        renderSelect(`misc.race`, `member${memberIndex}Race`, [
          { value: 'white', label: i18n.t('misc.race.white') },
          { value: 'latino', label: i18n.t('misc.race.latino') },
          { value: 'black', label: i18n.t('misc.race.black') },
          { value: 'asian', label: i18n.t('misc.race.asian') },
          { value: 'other', label: i18n.t('misc.other') }
        ]),
        renderSelect(`misc.relationship`, `member${memberIndex}Relationship`, [
          { value: 'child', label: i18n.t('misc.child') },
          { value: 'grandchild', label: i18n.t('misc.grandchild') },
          { value: 'spouse', label: i18n.t('misc.spouse') },
          { value: 'parent', label: i18n.t('misc.parent') },
          { value: 'grandparent', label: i18n.t('misc.grandparent') },
          { value: 'sibling', label: i18n.t('misc.sibling') },
          { value: 'friend', label: i18n.t('misc.friend') },
          { value: 'other', label: i18n.t('misc.other') }
        ]),
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
        m('a', { href: '#', onclick: (e) => { e.preventDefault(); i18n.setLanguage('en'); }, class: 'mr-4' }, 'English'),
        m('a', { href: '#', onclick: (e) => { e.preventDefault(); i18n.setLanguage('es'); } }, 'Español')
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
        renderSelect(`misc.race`, 'race', [
          { value: 'white', label: i18n.t('misc.race.white') },
          { value: 'latino', label: i18n.t('misc.race.latino') },
          { value: 'black', label: i18n.t('misc.race.black') },
          { value: 'asian', label: i18n.t('misc.race.asian') },
          { value: 'other', label: i18n.t('misc.other') }
        ]),
        renderSelect(`misc.primarylang`, 'primaryLanguage', [
          { value: 'english', label: i18n.t('misc.english') },
          { value: 'spanish', label: i18n.t('misc.spanish') },
          { value: 'other', label: i18n.t('misc.other') }
        ]),
        renderSelect(`misc.gender`, 'gender', [
          { value: 'male', label: i18n.t('misc.male') },
          { value: 'female', label: i18n.t('misc.female') },
          { value: 'optout', label: i18n.t('misc.prefernottosay') }
        ]),
        m('div', { class: 'mb-4' }, [
          m('label', { class: 'block text-gray-700 text-sm font-bold mb-2', for: 'dob' }, i18n.t('misc.dob')),
          m('div', { class: 'flex' }, [
            m('select', {
              class: `shadow appearance-none border ${errors['dobMonth'] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: 'dobMonth',
              oninput: function (ev) { household.head['dobMonth'] = ev.target.value; }
            }, [
              { value: '01', label: '01' },
              { value: '02', label: '02' },
              // Add more options for months
            ].map(option => m('option', { value: option.value }, option.label))),
            m('select', {
              class: `shadow appearance-none border ${errors['dobDay'] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: 'dobDay',
              oninput: function (ev) { household.head['dobDay'] = ev.target.value; }
            }, [
              { value: '01', label: '01' },
              { value: '02', label: '02' },
              // Add more options for days
            ].map(option => m('option', { value: option.value }, option.label))),
            m('select', {
              class: `shadow appearance-none border ${errors['dobYear'] ? 'border-red-500' : ''} rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline`,
              name: 'dobYear',
              oninput: function (ev) { household.head['dobYear'] = ev.target.value; }
            }, [
              { value: '2000', label: '2000' },
              { value: '2001', label: '2001' },
              // Add more options for years
            ].map(option => m('option', { value: option.value }, option.label)))
          ])
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
