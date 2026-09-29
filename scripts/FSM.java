import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.function.Predicate;

@FunctionalInterface
interface Action<S, E> {
    void execute(S from, E event, S to);
}

class Transition<S, E> {
    final S from;
    final E event;
    final S to;
    final Action<S, E> action;

    Transition(S from, E event, S to, Action<S, E> action) {
        this.from = from;
        this.event = event;
        this.to = to;
        this.action = action;
    }
}

public class FSM<S, E> {
    private S currentState;
    private final Map<S, Map<E, Transition<S, E>>> transitions = new HashMap<>();

    public FSM(S initialState) {
        this.currentState = Objects.requireNonNull(initialState, "初始状态不能为空");
    }

    /**
     * 支持链式调用）
     */
    public synchronized FSM<S, E> addTransition(S from, E event, S to, Action<S, E> action) {
        Map<E, Transition<S, E>> eventMap = transitions.computeIfAbsent(from, k -> new HashMap<>());
        if (eventMap.containsKey(event)) {
            throw new IllegalArgumentException("状态 " + from + " 针对事件 " + event + " 已存在转移规则！");
        }
        eventMap.put(event, new Transition<>(from, event, to, action));
        return this;
    }

    /**
     * 触发事件
     */
    public synchronized void fire(E event) {
        Map<E, Transition<S, E>> eventMap = transitions.get(currentState);
        if (eventMap == null) {
            throw new IllegalStateException("当前状态没有可用的转移规则：" + currentState);
        }

        Transition<S, E> transition = eventMap.get(event);
        if (transition == null) {
            throw new IllegalStateException("当前状态 " + currentState + " 不支持事件 " + event);
        }

        S oldState = currentState;
        S nextState = transition.to;

        // 先执行 Action，确保成功后再变更状态
        if (transition.action != null) {
            transition.action.execute(oldState, event, nextState);
        }

        // 状态更新
        this.currentState = nextState;
    }

    public synchronized S getCurrentState() {
        return currentState;
    }
}