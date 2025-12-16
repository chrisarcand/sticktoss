import { writable } from 'svelte/store';

// Auth store
export const user = writable(null);
export const isAuthenticated = writable(!!localStorage.getItem('token'));

// Skill level descriptions
export const skillLevels = {
  1: {
    label: 'Bender',
    description: "Can't stop without using the boards. Skating is labored and unstable—ankles rolling, choppy strides, limited mobility. Stickhandling is minimal; mostly just tries to get a piece of the puck. Shooting form is rough. Positional awareness is absent—chases the puck everywhere. Struggles with basic rules and gameplay concepts. Probably learning to skate as an adult or just started playing hockey."
  },
  2: {
    label: 'Pylon',
    description: "Can skate forward with moderate stability but crossovers are rough and backwards skating is shaky. Stops and starts without total control. Can receive and make simple passes but struggles under pressure. Shot exists but lacks consistency or power. Understands basic positions but gets caught out of place regularly. Easy to play around—opponents skate past like you're a stationary orange cone at practice. Might have been playing for years but the skills just haven't developed. Shows up and tries, which counts for something."
  },
  3: {
    label: 'Solid',
    description: "Competent skater in all directions with decent speed and agility. Can execute crossovers and transitions reasonably well. Handles the puck confidently in open ice and makes smart, tape-to-tape passes most of the time. Has a serviceable shot with occasional accuracy. Understands positioning and reads plays adequately. Reliable player who won't hurt the team but won't dominate either. May have played a bit growing up or has several years of dedicated beer league experience."
  },
  4: {
    label: 'Stud',
    description: "Strong, fluid skater with good edge work and quick acceleration. Handles the puck well in traffic and can make plays under pressure. Consistently accurate passer who sees the ice well. Has a legitimately dangerous shot—goalies respect it. Solid hockey IQ with good positioning and anticipation. Didn't grow up playing organized hockey but has excellent athleticism and years of experience that shows. Often one of the best players on the ice in most beer league games."
  },
  5: {
    label: 'Ringer',
    description: "Effortless, powerful skating with elite edge work and explosive speed. Stickhandles in a phone booth and protects the puck naturally. Makes high-difficulty passes look routine and sees plays developing before they happen. Can pick corners consistently and has a legitimately hard, accurate shot. Reads the game at a different speed than everyone else—always in the right position. Almost certainly played high school hockey, at least. Makes everyone else look slow. The guy who \"takes it easy\" and still dominates."
  }
};
