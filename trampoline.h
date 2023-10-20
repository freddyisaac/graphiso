#ifndef TRAMPOLINE_H
#define TRAMPOLINE_H

#include "rav1e.h"

typedef struct {
	RaConfig *rac;	
	RaContext *rax;
	RaFrame *raf;
} RAV1E;

extern RAV1E *new_rav1e();
extern void t_rav1e_config_default(RAV1E*);
extern void t_rav1e_clean(RAV1E*);
extern int t_rav1e_simple_setup(RAV1E*);
extern int t_rav1e_context_and_frame(RAV1E*);

extern int t_rav1e_simple_chromaticity(RAV1E*);

extern int t_rav1e_send(RAV1E*);

#endif // TRAMPOLINE_H

